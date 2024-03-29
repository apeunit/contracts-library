package aecl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/aeternity/aepp-contracts-library/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	// load postgresql driver
	_ "github.com/lib/pq"
)

// channels
var (
	db      *sql.DB
	proxies map[string]*httputil.ReverseProxy
	home    bytes.Buffer
)

// StartProxy starts the reverse proxy
func StartProxy(router *chi.Mux) (err error) {
	proxies = make(map[string]*httputil.ReverseProxy)
	// Prepare all the proxies
	for _, c := range Config.Compilers {
		compilerURL, err := url.Parse(c.URL)
		if err != nil {
			log.Printf("Error registerning proxy version %s at %s: %v", c.Version, c.URL, err)
		}
		if utils.IsEmptyStr(c.Version) {
			log.Printf("Error registerning proxy at %s: version field cannot be empty", c.URL)
		}
		log.Printf("Registered compiler %s at %s (Default: %v)", c.Version, c.URL, c.IsDefault)
		proxy := httputil.NewSingleHostReverseProxy(compilerURL)
		proxy.Transport = &LoggingTransport{}
		// TODO: ping the node before registering it
		// add to the proxies list
		proxies[c.Version] = proxy
		// if it is the default set it with empty string
		if c.IsDefault {
			proxies[""] = proxy
		}
	}

	// Open the db connection
	db, err = sql.Open("postgres", Config.DatabaseURL)
	db.SetMaxOpenConns(Config.Tuning.DbMaxOpenConns)
	db.SetMaxIdleConns(Config.Tuning.DbMaxIdleConns)
	if err != nil {
		log.Println("Error establishing connection to the database", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Println("Database ping failed", err)
		return
	}

	// setup the router
	// handle the home page
	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		data := struct {
			Version   string
			Compilers []CompilerSchema
			Header    string
		}{
			Version:   "1.0.0",
			Compilers: Config.Compilers,
			Header:    Config.Tuning.VersionHeaderName,
		}

		if req.Header.Get("Accept") == "application/json" {
			render.JSON(res, req, data)
			return
		}
		err = renderTemplate(Config.Web.HomeTemplatePath, data, res)
		return
	})
	// terms of service
	router.Get("/tos", func(res http.ResponseWriter, req *http.Request) {
		err = renderTemplate(Config.Web.TosTemplatePath, struct{}{}, res)
	})

	// handle static files
	log.Println("Serving assets from", Config.Web.AssetsFolderPath, "at", Config.Web.AssetsWebPath)
	fs := http.FileServer(http.Dir(Config.Web.AssetsFolderPath))
	router.Handle(Config.Web.AssetsWebPath, http.StripPrefix(path.Dir(Config.Web.AssetsWebPath), fs))
	// finally register the routes in the http module
	router.Handle("/*", http.HandlerFunc(HandleRequestAndRedirect))
	// return
	return
}

func renderTemplate(templatePath string, data interface{}, writer io.Writer) (err error) {
	// Build the homepage
	t, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		log.Println("Template build for", templatePath, " failed", err)
		return
	}
	// generate the page once and for all since we do
	// not support hot reloading of the configuration
	if err = t.Execute(writer, data); err != nil {
		log.Println("Template build for ", templatePath, " failed", err)
		return
	}
	return
}

// Contract the struct for peers
type Contract struct {
	Source       string                 `json:"code"`
	Options      map[string]interface{} `json:"options"`
	B2bH         string                 `json:"b2bh"`
	ResponseCode int                    `json:"response_code"`
	ResponseMsg  string                 `json:"response_msg"`
	Name         string                 `json:"name"`
}

func storeContract(contract *Contract) (err error) {
	// log.Printf("%#v", db.Stats())
	_, err = db.Exec(` INSERT INTO 
	contracts (id, source, response_code, response_msg) VALUES ( $1, $2, $3, $4) 
	ON CONFLICT (id) DO UPDATE SET 
	compilations = contracts.compilations + 1,
	updated_at = $5`,
		contract.B2bH,
		contract.Source,
		contract.ResponseCode,
		contract.ResponseMsg,
		time.Now())
	return
}

// HandleRequestAndRedirect Given a request send it to the appropriate url
func HandleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	rip := req.RemoteAddr
	path := req.URL.Path
	// get header request
	compilerVersion := req.Header.Get(Config.Tuning.VersionHeaderName)
	// resolve the request
	log.Println("Request from ", rip, " to ", path, " compiler version ", compilerVersion)
	proxy, found := proxies[compilerVersion]
	if !found {
		log.Println("Compiler version ", compilerVersion, " not found")
		return
	}
	// resolve the request
	proxy.ServeHTTP(res, req)
}

// LoggingTransport for roundtrips
type LoggingTransport struct {
	CapturedTransport http.RoundTripper
}

// RoundTrip function to save the request
func (t *LoggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// proxy ge requests
	if request.Method == http.MethodGet {
		response, err = http.DefaultTransport.RoundTrip(request)
		return
	}

	// proxy empty body requests
	if request.ContentLength <= 0 {
		response, err = http.DefaultTransport.RoundTrip(request)
		return
	}

	// proxy not compile requests
	if request.RequestURI != "/compile" {
		response, err = http.DefaultTransport.RoundTrip(request)
		return
	}
	// consume the request body buffer
	requestBody, err := ioutil.ReadAll(io.LimitReader(request.Body, Config.Tuning.RequestMaxBodySize))
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}
	// log the request time
	start := time.Now()
	// reset the buffer for the request
	request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	// execute the request
	response, err = http.DefaultTransport.RoundTrip(request)
	if err != nil {
		log.Println("Error response from backend", err)
		return
	}
	// store the contract
	contract := &Contract{}
	err = json.NewDecoder(bytes.NewReader(requestBody)).Decode(&contract)
	if err != nil {
		fmt.Println("Error decode contract", err)
		return
	}
	// compute the hash
	contract.B2bH = fmt.Sprintf("%x", blake2b.Sum256([]byte(contract.Source)))
	// get response data
	contract.ResponseCode = response.StatusCode
	// get response message
	if response.ContentLength > 0 {
		responseBody, _ := ioutil.ReadAll(io.LimitReader(response.Body, Config.Tuning.RequestMaxBodySize))
		contract.ResponseMsg = string(responseBody)
		// reset the buffer for the request
		response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
	}
	// store the contract
	go storeContract(contract)
	// print the reply
	log.Println("Response [", contract.ResponseCode, "] took ", time.Since(start))
	return
}

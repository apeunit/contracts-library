package aecl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/blake2b"
)

// channels
var (
	proxy *httputil.ReverseProxy
	db    *sql.DB
)

// StartProxy starts the reverse proxy
func StartProxy() (err error) {
	compilerURL, err := url.Parse(Config.CompilerURL)
	if err != nil {
		log.Println("Error starting the proxy ", err)
		return
	}
	log.Println("Proxy ready for ", compilerURL)
	proxy = httputil.NewSingleHostReverseProxy(compilerURL)
	proxy.Transport = &LoggingTransport{}

	// Open the db connection
	db, err = sql.Open("postgres", Config.DatabaseURL)
	if err != nil {
		log.Println("Error establishing connection to the database", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Println("Database ping failed", err)
		return
	}

	return
}

// Contract the struct for peers
type Contract struct {
	Source       string            `json:"code"`
	Options      map[string]string `json:"options"`
	B2bH         string            `json:"b2bh"`
	ResponseCode int               `json:"response_code"`
	ResponseMsg  string            `json:"response_msg"`
}

func storeContract(contract *Contract) (err error) {
	//
	_, err = db.Query(` INSERT INTO 
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
	log.Println("Request from ", rip, " to ", path)
	proxy.ServeHTTP(res, req)
}

type LoggingTransport struct {
	CapturedTransport http.RoundTripper
}

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
	requestBody, err := ioutil.ReadAll(io.LimitReader(request.Body, Config.Tuning.MaxBodySize))
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
		responseBody, _ := ioutil.ReadAll(io.LimitReader(response.Body, Config.Tuning.MaxBodySize))
		contract.ResponseMsg = string(responseBody)
		// reset the buffer for the request
		response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
	}
	// store the contract
	storeContract(contract)
	// print the reply
	log.Println("Response [", contract.ResponseCode, "] took ", time.Since(start))
	return
}

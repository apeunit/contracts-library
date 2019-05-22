package nodelb

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/mxmCherry/movavg"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"sort"
	"sync"
	"time"
)

// Cache variables
var (
	Nodes   *NodeCache
	Clients *bigcache.BigCache
	NodeMap map[string]Node
)

// Global Nodes variables
var (
	height uint64
)

// channels
var (
	tickChan chan time.Ticker
	doneChan chan bool
)

// NodeRank the ranked list of the nodes
type NodeRank []*Node

func (a NodeRank) Len() int      { return len(a) }
func (a NodeRank) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a NodeRank) Less(i, j int) bool {
	// NodeRnak should be sorted for
	// height descending
	// clients ascending
	// latency ascending
	// last contact ascending
	switch {
	case a[i].Height > a[j].Height:
		return true
	case a[i].Clients < a[j].Clients: // it would be better to deal with hits/s
		return true
	case a[i].Latency.Avg() < a[j].Latency.Avg():
		return true
	case a[i].LastContact.After(a[j].LastContact):
		return true
	default:
		return false
	}
}

// NodeCache keep the node sorted
type NodeCache struct {
	lock    sync.RWMutex     // read write lock
	NodeMap map[string]*Node // keep ipPort -> Rank
	Ranks   NodeRank         // ranked nodes
	Size    int              // size of the node cache
}

// GetNextNode get a node
func (nc *NodeCache) GetNextNode() (nodeIPPort string, err error) {
	nc.lock.RLock()
	defer nc.lock.RUnlock()
	if len(nc.Ranks) == 0 {
		err = fmt.Errorf("No nodes available")
		return
	}
	nodeIPPort = nc.Ranks[0].IPPort // get the first node since it is the one that is best suited to reply
	// reply to the node but update the client counters
	go func(nodeIPPort string) {
		nc.lock.Lock()
		defer nc.lock.Unlock()
		nc.NodeMap[nodeIPPort].Clients++
		nc.NodeMap[nodeIPPort].Hits++
		sort.Sort(nc.Ranks) // re-sort the ranks
	}(nodeIPPort)

	return
}

// Hit get a node
func (nc *NodeCache) Hit(nodeIPPort string) {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	nc.NodeMap[nodeIPPort].Hits++
}

// SetupCache setup the cache
func SetupCache() (err error) {
	// intialize the ticker
	// cache nodes and refresh before expiration
	Nodes = &NodeCache{
		NodeMap: make(map[string]*Node),
		Ranks:   NodeRank{},
		Size:    0,
		lock:    sync.RWMutex{},
	}
	if err != nil {
		return
	}
	// prepare the client cache
	cfg := bigcache.Config{
		LifeWindow: time.Duration(Config.ClientRetentionSec) * time.Second,
		OnRemoveWithReason: func(key string, val []byte, reason bigcache.RemoveReason) {
			// TODO: decrease client count for the node
			fmt.Println("Removing client ", key)
		},
	}
	// cache nodes and refresh before expiration
	Clients, err = bigcache.NewBigCache(cfg)
	if err != nil {
		return
	}
	return
}

// RegisterNode register a node
func RegisterNode(ipPort string) (nodeCount int) {
	Nodes.lock.Lock()
	defer Nodes.lock.Unlock()

	node, exists := Nodes.NodeMap[ipPort]
	if exists {
		fmt.Printf("Node (total: %d) %10s new: %5v height: %8d avg latency (sample %d):%f",
			Nodes.Size,
			node.IPPort,
			!exists,
			height,
			Config.Tuning.MovingAverageWindowSize,
			node.Latency.Avg())
		return Nodes.Size
	}

	// query the node height
	// TODO: probably unlock befor this call
	height, _, ttfb, err := TimeGet(fmt.Sprint("http://", ipPort, "/v2/key-blocks/current/height"))
	if err != nil {
		fmt.Println("Error registering node: ", err)
		return Nodes.Size
	}
	// new node
	node = &Node{
		IPPort:       ipPort,
		Height:       height,
		Latency:      movavg.NewSMA(Config.Tuning.MovingAverageWindowSize),
		RegisteredAt: time.Now(),
		LastContact:  time.Now(),
	}
	// registe the latency
	node.Latency.Add(ttfb.Seconds())
	// add the node
	Nodes.NodeMap[ipPort] = node            // add the node to the map
	Nodes.Ranks = append(Nodes.Ranks, node) // add the node to the ranks
	sort.Sort(Nodes.Ranks)                  // sort the ranks
	Nodes.Size++                            // increase the size

	// reply
	fmt.Printf("Node (total: %d) %10s new: %5v height: %8d avg latency (sample %d):%f", Nodes.Size, ipPort, !exists, height, Config.Tuning.MovingAverageWindowSize, node.Latency.Avg())
	// TODO: start scheduler
	return nodeCount
}

// GetNodeAddress get a node address
// there are 5 possibilities
// 1.  client is unkknowm                         => select the best available node
// 2.  client is known and node is ok             => use the same node
// 3.1 client is known and node is ko/out-of-sync => get a node with the same height (if available)
// 3.2 client is known and node is ko/out-of-sync => get a node with a height >= client expected height (if available)
// 3.3 client is known and node is ko/out-of-sync => get a node with the max(height) < client
func GetNodeAddress(clientIP string) (nodeIPPort string) {
	nodeIPB, err := Clients.Get(clientIP)
	if err != nil { // case 1
		nodeIPPort, _ = Nodes.GetNextNode()
		fmt.Printf("Client %s -> node: %s - cache err %v", clientIP, nodeIPPort, err)
		Clients.Set(clientIP, []byte(nodeIPPort))
		return
	}
	// TODO: verify node condition
	fmt.Printf("Client %s -> node: %s - from cache", clientIP, nodeIPPort)
	nodeIPPort = string(nodeIPB)
	go Nodes.Hit(nodeIPPort)
	return
}

// TimeGet make a request to an url and register the duration of the request
func TimeGet(url string) (height uint64, connectTime, ttfb time.Duration, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	var start, connect, dns, tlsHandshake time.Time
	trace := &httptrace.ClientTrace{

		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done: %v\n", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectTime = time.Since(connect)
			fmt.Printf("Connect time: %v\n", connectTime)
		},

		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
			fmt.Printf("Time from start to first byte: %v\n", ttfb)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	var myTransport http.RoundTripper = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: time.Second * 10,
	}

	resp, err := myTransport.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	type ChainHeight struct {
		Height uint64 `json:"height"`
	}
	defer resp.Body.Close()
	ch := ChainHeight{}
	json.NewDecoder(resp.Body).Decode(&ch)
	height = ch.Height
	fmt.Printf("Total time: %v\n", time.Since(start))
	return
}

// HandleRequestAndRedirect Given a request send it to the appropriate url
func HandleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	rip := req.RemoteAddr
	path := req.URL.Path
	nodeIP := GetNodeAddress(rip)
	url, _ := url.Parse(fmt.Sprint("http://", nodeIP, path))
	httputil.NewSingleHostReverseProxy(url).ServeHTTP(res, req)
}

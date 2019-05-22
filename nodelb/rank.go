package nodelb

import (
	"fmt"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/mxmCherry/movavg"
	"time"
)

// Node data relative to a aeternity node
type Node struct {
	IPPort       string      // the ip port of the node
	Latency      *movavg.SMA // moving average of the latencies of the node
	Height       uint64      // height of the node
	Clients      uint64      // number of connected client
	Hits         uint64      // number of hits to the node
	LastContact  time.Time   // last update
	RegisteredAt time.Time   // when the node was registered
	Ticker       time.Ticker // the ticker for updating the node info
	index        int         // The index of the item in the heap.
}

func (n Node) String() string {
	return fmt.Sprintf("Node %21s H:%10d, L:%f, C:%d, T:%s", n.IPPort, n.Height, n.Latency.Avg(), n.Clients, n.LastContact)
}

func NewNode(IPPort string) (node *Node) {
	return &Node{
		IPPort:  IPPort,
		Latency: movavg.NewSMA(Config.Tuning.MovingAverageWindowSize),
	}
}

// NodeManager manages the load balancer
type NodeManager struct {
	heap      *binaryheap.Heap
	maxHeight uint64
}

// NewManager create a new load Balance
func NewManager() *NodeManager {
	return &NodeManager{
		heap: binaryheap.NewWith(nodeRank),
	}
}

// Custom comparator (sort by IDs)
func nodeRank(a, b interface{}) int {

	nodeI := a.(*Node)
	nodeJ := b.(*Node)

	switch {
	case nodeI.Height < nodeJ.Height:
		return 1
	case nodeI.Height > nodeJ.Height:
		return -1
	case nodeI.Clients > nodeJ.Clients: // it would be better to deal with hits/s
		return 1
	case nodeI.Clients < nodeJ.Clients: // it would be better to deal with hits/s
		return -1
	case nodeI.Latency.Avg() > nodeJ.Latency.Avg():
		return 1
	case nodeI.Latency.Avg() < nodeJ.Latency.Avg():
		return -1
	case nodeI.LastContact.After(nodeJ.LastContact):
		return 1
	case nodeI.LastContact.Before(nodeJ.LastContact):
		return -1
	default:
		return 0
	}
}

// AddNode add node to the quee
func (nm *NodeManager) AddNode(node *Node) {
	// contact the node
	// register the height
	// register the latency
	node.Latency = movavg.NewSMA(Config.Tuning.MovingAverageWindowSize)
	fmt.Println("Add node ", node.IPPort)
	nm.heap.Push(node)
}

// Dump dump the content of the heap
func (nm *NodeManager) Dump() {
	i := nm.heap.Iterator()
	for i.Next() {
		n := i.Value().(*Node)
		ids := i.Index()
		fmt.Printf("[%d] %s\n", ids, n)
	}
}

// GetBestNode retrieve the best node
func (nm *NodeManager) GetBestNode() (n *Node) {
	// remove the head
	v, _ := nm.heap.Pop()
	n = v.(*Node)
	// increase the number of clients
	n.Clients++
	// reinsert the node
	nm.heap.Push(n)
	return
}

// Step 1. register a node
// Step 2.

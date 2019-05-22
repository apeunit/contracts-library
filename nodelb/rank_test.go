package nodelb

import (
	"fmt"
	"testing"
)

func TestNodePriorityQueue_AddNode(t *testing.T) {
	Config = ConfigSchema{}
	Config.Defaults()
	// this is the queue
	pq := NewManager()

	tests := []struct {
		name string
		node *Node
	}{
		{"2", &Node{IPPort: "99.136.37.63:3013", Height: 100, Clients: 2}},
		{"2", &Node{IPPort: "35.178.61.73:3013", Height: 20}},
		{"2", &Node{IPPort: "35.177.192.219:3013", Height: 50}},
		{"2", &Node{IPPort: "18.136.37.63:3013", Height: 100}},
		{"2", &Node{IPPort: "52.220.198.72:3013", Height: 10}},
		{"2", &Node{IPPort: "88.220.198.72:3013", Height: 50, Clients: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq.AddNode(tt.node)
		})
	}

	pq.Dump()

	head := pq.GetBestNode()
	if head.IPPort != "18.136.37.63:3013" {
		t.Logf("expected: %s got %s", "18.136.37.63:3013", head.IPPort)
		t.Fail()
	}
	fmt.Println("Round 2")
	head = pq.GetBestNode()
	pq.Dump()
	fmt.Println("Round 3")
	head = pq.GetBestNode()
	if head.IPPort != "99.136.37.63:3013" {
		t.Logf("expected: %s got %s", "99.136.37.63:3013", head.IPPort)
		t.Fail()
	}
	pq.Dump()

	t.Fatalf("Finito")
}

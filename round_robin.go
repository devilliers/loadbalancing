package loadbalancing

import (
	"errors"
	"math/rand"
	"time"
)

// RoundRobin is
type RoundRobin struct {
	current    int
	bestMethod string
}

// TODO: implement a system for choosing the bestMethod

// Execute chooses the balancing method based on the current bestMethod string
func (r *RoundRobin) Execute(nodes []*Node) *Node {
	// var node *Node
	node := &Node{}
	switch r.bestMethod {
	case "weightedElect":
		node, _ = r.weightedElect(nodes)
	case "leastConnElect":
		node, _ = r.leastConnElect(nodes)
	case "weightedLeastConnElect":
		node, _ = r.weightedLeastConnElect(nodes)
	case "randElect":
		node, _ = r.randElect(nodes)
	default:
		node, _ = r.standardElect(nodes)
	}
	return node
}

// StandardElect is a non-weighted standard round_robin elect method: goes to next node
func (r *RoundRobin) standardElect(nodes []*Node) (*Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("no nodes to which to balance traffic")
	}

	if r.current >= len(nodes) {
		r.current = 0
	}

	r.current++

	return nodes[r.current], nil
}

// WeightedElect is a weighted round robin
func (r *RoundRobin) weightedElect(nodes []*Node, incr ...float32) (*Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("no nodes")
	}

	if len(incr) == 0 {
		incr = []float32{1.0}
	}

	if r.current >= len(nodes) {
		r.current = 0
	}

	node := nodes[r.current]

	// weightStore := nodes[r.current].weight
	if node.currentWeight <= incr[0] {
		// go to next node, reset current
		r.current++
		node.currentWeight = node.originalWeight
	} else {
		node.currentWeight -= incr[0]
	}

	return nodes[r.current], nil
}

// LeastConnElect is a least-connections based election method
func (r *RoundRobin) leastConnElect(nodes []*Node) (*Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("no nodes")
	}

	leastConnectionsNode := nodes[0]
	for i, node := range nodes {
		if node.currentConnections < leastConnectionsNode.currentConnections {
			leastConnectionsNode = nodes[i]
		}
	}

	return leastConnectionsNode, nil

}

func (r *RoundRobin) weightedLeastConnElect(nodes []*Node) (*Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("no nodes")
	}

	for _, node := range nodes {
		node.combinedWeight = float32(node.currentConnections) / node.originalWeight
	}

	min := nodes[0]
	for _, e := range nodes {
		if e.combinedWeight < min.combinedWeight {
			min = e
		}
	}
	return min, nil
}

func (r *RoundRobin) randElect(nodes []*Node) (*Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("no nodes")
	}
	rand.Seed(time.Now().UnixNano())
	return nodes[rand.Intn(len(nodes))], nil
}

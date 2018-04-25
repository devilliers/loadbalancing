package loadbalancing

// Node is an element to which the loadbalancer balances
type Node struct {
	originalWeight     float32
	currentWeight      float32
	currentConnections int
	combinedWeight     float32
}

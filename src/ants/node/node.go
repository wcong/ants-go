package node

type Node struct {
	Ip   string
	Port int
}

func NewNode(ip string, port int) *Node {
	return &Node{ip, port}
}

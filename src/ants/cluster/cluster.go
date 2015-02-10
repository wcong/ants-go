package cluster

import (
	. "ants/node"
)

type Cluster struct {
	name       string
	nodeList   []*Node
	masterNode *Node
}

func NewCluster(name string) *Cluster {
	return &Cluster{name, make([]*Node, 0), nil}
}

func (this *Cluster) SetName(name string) {
	this.name = name
}
func (this *Cluster) GetName() string {
	return this.name
}
func (this *Cluster) GetNodeList() []*Node {
	return this.nodeList
}

func (this *Cluster) AddNode(node *Node) {
	if this.nodeList == nil {
		this.nodeList = make([]*Node, 3)
	}
	append(this.nodeList, node)
}

func (this *Cluster) GetMasterNode() *Node {
	return this.masterNode
}

func (this *Cluster) SetMasterNode(node *Node) {
	this.masterNode = node
}

package node

import (
	"ants/conf"
)

type Cluster struct {
	Name       string
	NodeList   []*Node
	LocalNode  *Node
	MasterNode *Node
}

func NewCluster(settings *conf.Settings, localNode *Node) *Cluster {
	return &Cluster{settings.Name, make([]*Node, 0), localNode, localNode}
}

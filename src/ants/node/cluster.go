package node

import (
	"ants/conf"
)

type Cluster struct {
	Name       string
	NodeList   []*NodeInfo
	LocalNode  *NodeInfo
	MasterNode *NodeInfo
}

func NewCluster(settings *conf.Settings, localNode *NodeInfo) *Cluster {
	return &Cluster{settings.Name, make([]*NodeInfo, 0), localNode, localNode}
}
func (this *Cluster) AddNode(nodeInfo *NodeInfo) {
	this.NodeList = append(this.NodeList, nodeInfo)
	if this.LocalNode == this.MasterNode {
		this.ElectMaster()
	}
}
func (this *Cluster) ElectMaster() {

}

package node

import (
	"ants/util"
)

type Cluster struct {
	Name       string
	NodeList   []*NodeInfo
	LocalNode  *NodeInfo
	MasterNode *NodeInfo
}

func NewCluster(settings *util.Settings, localNode *NodeInfo) *Cluster {
	cluster := &Cluster{settings.Name, make([]*NodeInfo, 0), localNode, localNode}
	cluster.NodeList = append(cluster.NodeList, localNode)
	return cluster
}
func (this *Cluster) AddNode(nodeInfo *NodeInfo) {
	this.NodeList = append(this.NodeList, nodeInfo)
	if this.LocalNode == this.MasterNode {
		this.ElectMaster()
	}
}
func (this *Cluster) ElectMaster() *NodeInfo {
	this.MasterNode = this.NodeList[0]
	return this.MasterNode
}

package node

import (
	"log"
)

/*
* this is rpc method of cluster
**/
func (this *RPCer) letMeIn(ip string, port int) error {
	client, err := this.dial(ip, port)
	if err != nil {
		return err
	}
	request := new(LeftMeInRequest)
	request.NodeInfo = this.node.NodeInfo
	response := new(LeftMeInResponse)
	err = client.Call("RPCer.LetMeIn", request, response)
	if response.Result {
		this.connMap[response.NodeInfo.Name] = client
		this.node.AddNodeToCluster(response.NodeInfo)
		this.node.Cluster.MakeMasterNode(response.NodeInfo.Name)
	} else {
		client.Close()
		this.letMeIn(response.NodeInfo.Ip, response.NodeInfo.Port)
	}
	return err
}

// handle the join request,
// if this is the master node ,let it join and connect to server
// else send it  master node ,let it talk to master
func (this *RPCer) LetMeIn(request *LeftMeInRequest, response *LeftMeInResponse) error {
	if this.node.Cluster.IsMasterNode() {
		this.node.Join()
		response.Result = true
		response.NodeInfo = this.node.NodeInfo
		this.connect(request.NodeInfo)
		this.node.Ready()
	} else {
		response.Result = false
		response.NodeInfo = this.node.Cluster.GetMasterNode()
	}
	return nil
}

// for now it is for master connect to slave
func (this *RPCer) connect(nodeInfo *NodeInfo) error {
	client, err := this.dial(nodeInfo.Ip, nodeInfo.Port)
	if err != nil {
		return err
	}
	request := new(RpcBase)
	response := new(RpcBase)
	err = client.Call("RPCer.Connect", request, response)
	if err == nil {
		this.connMap[response.NodeInfo.Name] = client
		this.node.AddNodeToCluster(response.NodeInfo)
	}
	return err
}

// just tell who i am
func (this *RPCer) Connect(request *RpcBase, response *RpcBase) error {
	response.Result = true
	response.NodeInfo = this.node.NodeInfo
	return nil
}

// top node
func (this *RPCer) stopNode(nodeName string) error {
	stopRequest := &StopRequest{}
	stopRequest.NodeInfo = this.node.NodeInfo
	stopResponse := &StopResponse{}
	err := this.connMap[nodeName].Call("RPCer.StopNode", stopRequest, stopResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (this *RPCer) StopNode(request *StopRequest, response *StopResponse) error {
	this.node.StopCrawl()
	return nil
}

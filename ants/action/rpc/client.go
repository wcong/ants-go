package rpc

import (
	"github.com/wcong/ants-go/ants/action"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type RpcClient struct {
	node    *node.Node
	connMap map[string]*rpc.Client
}

func NewRpcClient(node *node.Node) *RpcClient {
	connMap := make(map[string]*rpc.Client)
	return &RpcClient{node, connMap}
}

// start a rpc client
// if ok , store if
func (this *RpcClient) Dial(ip string, port int) (*rpc.Client, error) {
	client, err := jsonrpc.Dial(RPC_TYPE, ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
	}
	return client, err
}

/*
* this is rpc method of cluster
**/
func (this *RpcClient) LetMeIn(ip string, port int) error {
	client, err := this.Dial(ip, port)
	if err != nil {
		return err
	}
	request := new(action.LeftMeInRequest)
	request.NodeInfo = this.node.NodeInfo
	response := new(action.LeftMeInResponse)
	err = client.Call("RpcServer.LetMeIn", request, response)
	if response.Result {
		this.connMap[response.NodeInfo.Name] = client
		this.node.AddNodeToCluster(response.NodeInfo)
		this.node.Cluster.MakeMasterNode(response.NodeInfo.Name)
	} else {
		client.Close()
		this.LetMeIn(response.NodeInfo.Ip, response.NodeInfo.Port)
	}
	return err
}

// for now it is for master connect to slave
func (this *RpcClient) Connect(ip string, port int) error {
	client, err := this.Dial(ip, port)
	if err != nil {
		return err
	}
	request := new(action.RpcBase)
	response := new(action.RpcBase)
	err = client.Call("RpcServer.Connect", request, response)
	if err == nil {
		this.connMap[response.NodeInfo.Name] = client
		this.node.AddNodeToCluster(response.NodeInfo)
	}
	return err
}

// top node
func (this *RpcClient) StopNode(nodeName string) error {
	stopRequest := &action.StopRequest{}
	stopRequest.NodeInfo = this.node.NodeInfo
	stopResponse := &action.StopResponse{}
	err := this.connMap[nodeName].Call("RpcServer.StopNode", stopRequest, stopResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

/*
*	this is rpc method of request
**/
func (this *RpcClient) Distribute(nodeName string, request *http.Request) error {
	distributeRequest := &action.DistributeRequest{}
	distributeRequest.NodeInfo = this.node.NodeInfo
	distributeRequest.Request = request
	distributeReqponse := &action.DistributeReqponse{}
	err := this.connMap[nodeName].Call("RpcServer.AcceptRequest", distributeRequest, distributeReqponse)
	if err != nil {
		log.Println(err)
	} else {
		this.node.AddToCrawlingQuene(request)
	}
	return err
}

func (this *RpcClient) StartSpider(nodeName, spiderName string) error {
	startRequest := &action.RpcBase{
		NodeInfo: this.node.NodeInfo,
	}
	startResponse := &action.RpcBase{}
	err := this.connMap[nodeName].Call("RpcServer.StartSpider", startRequest, startResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

// for slave send crawl result to master
func (this *RpcClient) ReportResult(nodeName string, result *crawler.ScrapeResult) error {
	reportRequest := &action.ReportRequest{}
	reportRequest.NodeInfo = this.node.NodeInfo
	reportRequest.ScrapeResult = result
	reportResponse := &action.ReportResponse{}
	err := this.connMap[nodeName].Call("RpcServer.AcceptResult", reportRequest, reportResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

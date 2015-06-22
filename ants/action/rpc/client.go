package rpc

import (
	"github.com/wcong/ants-go/ants/action"
	"github.com/wcong/ants-go/ants/cluster"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

type RpcClient struct {
	node    node.Node
	cluster cluster.Cluster
	connMap map[string]*rpc.Client
}

func NewRpcClient(node node.Node, cluster cluster.Cluster) *RpcClient {
	connMap := make(map[string]*rpc.Client)
	return &RpcClient{node, cluster, connMap}
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
	request.NodeInfo = this.node.GetNodeInfo()
	response := new(action.LeftMeInResponse)
	err = client.Call("RpcServer.LetMeIn", request, response)
	if response.Result {
		this.connMap[response.NodeInfo.Name] = client
		this.cluster.AddNode(response.NodeInfo)
		this.cluster.MakeMasterNode(response.NodeInfo.Name)
	} else {
		client.Close()
		this.LetMeIn(response.NodeInfo.Ip, response.NodeInfo.Port)
	}
	return err
}
func (this *RpcClient) Start() {
	go func() {
		for {
			this.Detect()
			time.Sleep(2 * time.Second)
		}
	}()
}

// judge if node is down
func (this *RpcClient) Detect() {
	request := new(action.RpcBase)
	response := new(action.RpcBase)
	for key, conn := range this.connMap {
		err := conn.Call("RpcServer.IsAlive", request, response)
		if err != nil {
			log.Println(err)
			log.Println("node:", key, "is dead ,so remove it")
			delete(this.connMap, key)
			if this.node.IsMasterNode() {
				this.cluster.DeleteDeadNode(key)
			}
		}
	}
}

func (this *RpcClient) SyncClusterInfo() {

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
		this.cluster.AddNode(response.NodeInfo)
	}
	return err
}

// top node
func (this *RpcClient) StopNode(nodeName string) error {
	stopRequest := &action.StopRequest{}
	stopRequest.NodeInfo = this.node.GetNodeInfo()
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
	distributeRequest.NodeInfo = this.node.GetNodeInfo()
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
		NodeInfo: this.node.GetNodeInfo(),
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
	reportRequest.NodeInfo = this.node.GetNodeInfo()
	reportRequest.ScrapeResult = result
	reportResponse := &action.ReportResponse{}
	err := this.connMap[nodeName].Call("RpcServer.AcceptResult", reportRequest, reportResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (this RpcClient) CloseSpider(nodeName, spiderName string) error {
	closeRequest := &action.CloseSpiderRequest{}
	closeRequest.NodeInfo = this.node.GetNodeInfo()
	closeRequest.SpiderName = spiderName
	closeResponse := &action.CloseSpiderResponse{}
	err := this.connMap[nodeName].Call("RpcServer.CloseSpider", closeRequest, closeResponse)
	if err != nil {
		log.Panicln(err)
	}
	return err
}

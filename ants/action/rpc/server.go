package rpc

import (
	"github.com/wcong/ants-go/ants/action"
	"github.com/wcong/ants-go/ants/cluster"
	"github.com/wcong/ants-go/ants/node"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

const (
	RPC_TYPE = "tcp"
)

type RpcServer struct {
	node        node.Node
	cluster     cluster.Cluster
	port        int
	rpcClient   action.RpcClientAnts
	reporter    action.Watcher
	distributer action.Watcher
}

func NewRpcServer(node node.Node, cluster cluster.Cluster, port int, rpcClient action.RpcClientAnts, reporter, distributer action.Watcher) *RpcServer {
	rpcServer := &RpcServer{node, cluster, port, rpcClient, reporter, distributer}
	rpcServer.start()
	return rpcServer
}

// start a rpc server
func (this *RpcServer) server() {
	rpc.Register(this)
	listener, e := net.Listen(RPC_TYPE, ":"+strconv.Itoa(this.port))
	if e != nil {
		log.Println(e)
		return
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Println(err)
		} else {
			log.Println("new connection")
			go jsonrpc.ServeConn(conn)
		}
	}
}

// start a dead loop for server
func (this *RpcServer) start() {
	go this.server()
}

// is server alive
func (this *RpcServer) IsAlive(request *action.RpcBase, response *action.RpcBase) error {
	return nil
}

// handle the join request,
// if this is the master node ,let it join and connect to server
// else send it  master node ,let it talk to master
func (this *RpcServer) LetMeIn(request *action.LeftMeInRequest, response *action.LeftMeInResponse) error {
	if this.node.IsMasterNode() {
		this.node.Join()
		response.Result = true
		response.NodeInfo = this.node.GetNodeInfo()
		this.rpcClient.Connect(request.NodeInfo.Ip, request.NodeInfo.Port)
		this.node.Ready()
	} else {
		response.Result = false
		response.NodeInfo = this.cluster.GetMasterNode()
	}
	return nil
}

// just tell who i am
func (this *RpcServer) Connect(request *action.RpcBase, response *action.RpcBase) error {
	response.Result = true
	response.NodeInfo = this.node.GetNodeInfo()
	return nil
}

func (this *RpcServer) StopNode(request *action.StopRequest, response *action.StopResponse) error {
	this.stopNode()
	return nil
}

// expose method ,for accept method
func (this *RpcServer) AcceptRequest(request *action.DistributeRequest, response *action.DistributeReqponse) error {
	this.node.AcceptRequest(request.Request)
	return nil
}

// start spider for slave
// for now just start reporter
func (this *RpcServer) StartSpider(request *action.DistributeRequest, response *action.DistributeReqponse) error {
	this.reporter.Start()
	return nil
}

// for master accept crawl result
// if no more request to crawl stop it
// notice:
// *		local request do not go through this way
// *		close action also start by reporter
func (this *RpcServer) AcceptResult(request *action.ReportRequest, response *action.ReportResponse) error {
	this.node.AcceptResult(request.ScrapeResult)
	spiderName := request.ScrapeResult.Request.SpiderName
	if this.node.CanWeStopSpider(spiderName) {
		for _, nodeInfo := range this.cluster.GetAllNode() {
			if this.node.IsMe(nodeInfo.Name) {
				this.node.CloseSpider(spiderName)
			} else {
				this.rpcClient.CloseSpider(nodeInfo.Name, spiderName)
			}
		}
	}
	if this.node.IsStop() {
		for _, nodeInfo := range this.cluster.GetAllNode() {
			if this.node.IsMe(nodeInfo.Name) {
				this.stopNode()
			} else {
				this.rpcClient.StopNode(nodeInfo.Name)
			}
		}
	}
	return nil
}

// stop the node
// crawler ,reporter and distributer
func (this *RpcServer) stopNode() {
	this.node.StopCrawl()
	this.reporter.Stop()
	this.distributer.Stop()
}

func (this *RpcServer) CloseSpider(request *action.CloseSpiderRequest, response *action.CloseSpiderResponse) error {
	this.node.CloseSpider(request.SpiderName)
	return nil
}

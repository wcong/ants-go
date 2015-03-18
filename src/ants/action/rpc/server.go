package rpc

import (
	"ants/action"
	"ants/node"
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
	node        *node.Node
	port        int
	rpcClient   action.RpcClientAnts
	reporter    action.Watcher
	distributer action.Watcher
}

func NewRpcServer(node *node.Node, port int, rpcClient action.RpcClientAnts, reporter, distributer action.Watcher) *RpcServer {
	rpcServer := &RpcServer{node, port, rpcClient, reporter, distributer}
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

// handle the join request,
// if this is the master node ,let it join and connect to server
// else send it  master node ,let it talk to master
func (this *RpcServer) LetMeIn(request *action.LeftMeInRequest, response *action.LeftMeInResponse) error {
	if this.node.IsMasterNode() {
		this.node.Join()
		response.Result = true
		response.NodeInfo = this.node.NodeInfo
		this.rpcClient.Connect(request.NodeInfo.Ip, request.NodeInfo.Port)
		this.node.Ready()
	} else {
		response.Result = false
		response.NodeInfo = this.node.Cluster.GetMasterNode()
	}
	return nil
}

// just tell who i am
func (this *RpcServer) Connect(request *action.RpcBase, response *action.RpcBase) error {
	response.Result = true
	response.NodeInfo = this.node.NodeInfo
	return nil
}

func (this *RpcServer) StopNode(request *action.StopRequest, response *action.StopResponse) error {
	this.node.StopCrawl()
	return nil
}

// expose method ,for accept method
func (this *RpcServer) AcceptRequest(request *action.DistributeRequest, response *action.DistributeReqponse) error {
	this.node.AcceptRequest(request.Request)
	return nil
}

// for master accept crawl result
func (this *RpcServer) AcceptResult(request *action.ReportRequest, response *action.ReportResponse) error {
	this.node.AcceptResult(request.ScrapeResult)
	if this.node.IsStop() {
		for _, nodeInfo := range this.node.GetAllNodeForClose() {
			if this.node.IsMe(nodeInfo.Name) {
				this.node.StopCrawl()
			} else {
				this.rpcClient.StopNode(nodeInfo.Name)
			}
		}
	}
	return nil
}

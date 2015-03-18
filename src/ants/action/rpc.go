package action

import (
	"ants/crawler"
	"ants/http"
	"net/rpc"
)

/*
*	ants rpc server and client should server those function
*
**/

/*
*	RpcServer
**/
type RpcServer interface {
}

/*
*	RpcServerCrawl
**/
type RpcServerCrawl interface {
	AcceptRequest(request *DistributeRequest, response *DistributeReqponse) error
	AcceptResult(request *ReportRequest, response *ReportResponse) error
}

/*
*	RpcServerCluster
**/
type RpcServerCluster interface {
	LetMeIn(request *LeftMeInRequest, response *LeftMeInResponse) error
	Connect(request *RpcBase, response *RpcBase) error
	StopNode(request *StopRequest, response *StopResponse) error
}

/*
*	RpcServerAnts
**/
type RpcServerAnts interface {
	RpcServer
	RpcServerCrawl
	RpcServerCluster
}

/*
*	RpcClient
*
**/
type RpcClient interface {
	Dial(ip string, port int) (*rpc.Client, error)
}

/*
*	RpcClientCluster
**/
type RpcClientCluster interface {
	LetMeIn(ip string, port int) error
	Connect(ip string, port int) error
}

/*
*	RpcClientCrawl
**/
type RpcClientCrawl interface {
	StopNode(nodeName string) error
	Distribute(nodeName string, request *http.Request) error
	ReportResult(nodeName string, result *crawler.ScrapeResult) error
}

/*
*	RpcClientAnts
**/
type RpcClientAnts interface {
	RpcClient
	RpcClientCrawl
	RpcClientCluster
}

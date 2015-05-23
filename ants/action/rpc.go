package action

import (
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
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
	IsAlive(request *RpcBase, response *RpcBase) error
}

/*
*	RpcServerCrawl
**/
type RpcServerCrawl interface {
	AcceptRequest(request *DistributeRequest, response *DistributeReqponse) error
	AcceptResult(request *ReportRequest, response *ReportResponse) error
	CloseSpider(request *CloseSpiderRequest, response *CloseSpiderResponse) error
}

/*
*	RpcServerCluster
**/
type RpcServerCluster interface {
	LetMeIn(request *LeftMeInRequest, response *LeftMeInResponse) error
	Connect(request *RpcBase, response *RpcBase) error
	StopNode(request *StopRequest, response *StopResponse) error
	StartSpider(request *DistributeRequest, response *DistributeReqponse) error
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
	Detect() // if client connection is down
	Start()  // start dead loop for detect
}

/*
*	RpcClientCluster
**/
type RpcClientCluster interface {
	LetMeIn(ip string, port int) error
	Connect(ip string, port int) error
	StartSpider(nodeName, spiderName string) error
}

/*
*	RpcClientCrawl
**/
type RpcClientCrawl interface {
	StopNode(nodeName string) error
	Distribute(nodeName string, request *http.Request) error
	ReportResult(nodeName string, result *crawler.ScrapeResult) error
	CloseSpider(nodeName, spiderName string) error
}

/*
*	RpcClientAnts
**/
type RpcClientAnts interface {
	RpcClient
	RpcClientCrawl
	RpcClientCluster
}

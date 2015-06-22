package watcher

import (
	"github.com/wcong/ants-go/ants/action"
	"github.com/wcong/ants-go/ants/cluster"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"log"
	"time"
)

/*
what a Distributer do
*	status,running|parse|stop
*	distribute a request,by some strategy
*	*
*/
const (
	DISTRIBUTE_RUNING = iota
	DISTRIBUTE_PAUSE
	DISTRIBUTE_STOP
	DISTRIBUTE_STOPED
)

type Distributer struct {
	Status    int
	Cluster   cluster.Cluster
	Node      node.Node
	LastIndex int
	RpcClient action.RpcClientAnts
}

func NewDistributer(node node.Node, cluster cluster.Cluster, rpcClient action.RpcClientAnts) *Distributer {
	return &Distributer{DISTRIBUTE_STOPED, cluster, node, 0, rpcClient}
}

func (this *Distributer) IsStop() bool {
	return this.Status == DISTRIBUTE_STOPED
}

func (this *Distributer) IsPause() bool {
	return this.Status == DISTRIBUTE_PAUSE
}
func (this *Distributer) Pause() {
	if this.Status == DISTRIBUTE_RUNING {
		this.Status = DISTRIBUTE_PAUSE
	}
}
func (this *Distributer) Unpause() {
	if this.Status == DISTRIBUTE_PAUSE {
		this.Status = DISTRIBUTE_RUNING
	}
}

func (this *Distributer) Stop() {
	this.Status = DISTRIBUTE_STOP
}

func (this *Distributer) Start() {
	if this.Status == DISTRIBUTE_RUNING {
		return
	}
	for {
		if this.IsStop() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	this.Status = DISTRIBUTE_RUNING
	this.Run()
}

// dead loop cluster pop request
func (this *Distributer) Run() {
	log.Println("start distributer")
	for {
		if this.Status == DISTRIBUTE_STOP {
			this.Status = DISTRIBUTE_STOPED
			break
		}
		if this.IsPause() {
			time.Sleep(1 * time.Second)
			continue
		}
		request := this.Cluster.PopRequest()
		if request == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		this.Distribute(request)
		log.Println(request.SpiderName, ":distribute:", request.NodeName, ":request:", request.GoRequest.URL.String())
		if this.Node.IsMe(request.NodeName) {
			this.Node.DistributeRequest(request)
		} else {
			this.RpcClient.Distribute(request.NodeName, request)
		}
	}
	log.Println("stop distributer")
}

// if cookiejar > 0 means it require cookie context ,so we should send it to where it come from
// else distribute it by order
func (this *Distributer) Distribute(request *http.Request) {
	if request.CookieJar > 0 {
		return
	} else {
		nodeList := this.Cluster.GetAllNode()
		if this.LastIndex >= len(nodeList) {
			this.LastIndex = 0
		}
		nodeName := nodeList[this.LastIndex].Name
		request.NodeName = nodeName
		this.LastIndex += 1
	}
}

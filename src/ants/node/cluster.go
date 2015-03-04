package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"log"
	"sync"
)

/*
what a cluster would do
* 	init by local node
*	add a node
*	choose a master node
*	distribute a request
* 	accept crawl result
*	add a request
*/

// cluster status
// *		init:where every thing have init
// *		join:try to connect to other node ,if not make itself master,else ,get other master
// *		election(option):when circle is builded a start to elect a master
// * 	ready:ready to start crawl
const (
	CLUSTER_STATUS_INIT = iota
	CLUSTER_STATUS_JOIN
	CLUSTER_STATUS_ELECTION
	CLUSTER_STATUS_READY
)

// basic cluster infomation
type ClusterInfo struct {
	Status     int
	Name       string
	NodeList   []*NodeInfo
	LocalNode  *NodeInfo
	MasterNode *NodeInfo
}

// receive basic request and record crawled requets
type RequestStatus struct {
	CrawledMap   map[string]int
	CrawlingMap  map[string]map[string]*http.Request
	WaitingQuene *crawler.RequestQuene
}

type Cluster struct {
	ClusterInfo   *ClusterInfo
	RequestStatus *RequestStatus
	mutex         *sync.Mutex
}

func NewCluster(settings *util.Settings, localNode *NodeInfo) *Cluster {
	clusterInfo := &ClusterInfo{CLUSTER_STATUS_INIT, settings.Name, make([]*NodeInfo, 0), localNode, nil}
	requestStatus := &RequestStatus{}
	requestStatus.CrawledMap = make(map[string]int)
	requestStatus.CrawlingMap = make(map[string]map[string]*http.Request)
	requestStatus.WaitingQuene = crawler.NewRequestQuene()
	cluster := &Cluster{clusterInfo, requestStatus, new(sync.Mutex)}
	cluster.ClusterInfo.NodeList = append(cluster.ClusterInfo.NodeList, localNode)
	return cluster
}

// is local node master node
func (this *Cluster) IsMasterNode() bool {
	if this.ClusterInfo.MasterNode == nil {
		return false
	}
	return this.ClusterInfo.LocalNode.Name == this.ClusterInfo.MasterNode.Name
}

// add a node to cluster node list
func (this *Cluster) AddNode(nodeInfo *NodeInfo) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, node := range this.ClusterInfo.NodeList {
		if node.Name == nodeInfo.Name {
			return
		}
	}
	this.ClusterInfo.NodeList = append(this.ClusterInfo.NodeList, nodeInfo)
}

// make master node by node name
func (this *Cluster) MakeMasterNode(nodeName string) {
	for _, node := range this.ClusterInfo.NodeList {
		if node.Name == nodeName {
			this.ClusterInfo.MasterNode = node
		}
	}
}

// choose a new master node
func (this *Cluster) ElectMaster() *NodeInfo {
	this.ClusterInfo.MasterNode = this.ClusterInfo.NodeList[0]
	for _, node := range this.ClusterInfo.NodeList {
		if this.ClusterInfo.MasterNode.Name > node.Name {
			this.ClusterInfo.MasterNode = node
		}
	}
	return this.ClusterInfo.MasterNode
}

// pop a request from waiting quene
// add to crawling quenu
func (this *Cluster) PopRequest() *http.Request {
	return this.RequestStatus.WaitingQuene.Pop()
}

// record distribute request job
func (this *Cluster) AddToCrawlingQuene(request *http.Request) {
	requestMap, ok := this.RequestStatus.CrawlingMap[request.NodeName]
	if !ok {
		requestMap = make(map[string]*http.Request)
		this.RequestStatus.CrawlingMap[request.NodeName] = requestMap
	}
	requestMap[request.UniqueName] = request
}

// a request job is done
// delete it from crawling quene
// add crawled num
func (this *Cluster) Crawled(nodeName, requestHashName string) {
	requestMap, ok := this.RequestStatus.CrawlingMap[nodeName]
	if !ok {
		log.Println("none node :" + nodeName)
		return
	}
	delete(requestMap, requestHashName)
}

// add a request to quene
func (this *Cluster) AddRequest(request *http.Request) {
	this.RequestStatus.WaitingQuene.Push(request)
}

// get master node name
func (this *Cluster) GetMasterName() string {
	return this.ClusterInfo.MasterNode.Name
}

// is all loop stop
func (this *Cluster) IsStop() bool {
	if !this.RequestStatus.WaitingQuene.IsEmpty() {
		return false
	}
	for _, requestMap := range this.RequestStatus.CrawlingMap {
		if len(requestMap) > 0 {
			return false
		}
	}
	return true
}

func (this *Cluster) HasNode(nodeName string) bool {
	for _, node := range this.ClusterInfo.NodeList {
		if node.Name == nodeName {
			return true
		}
	}
	return false
}

// is cluster ready for crawl
func (this *Cluster) IsReady() bool {
	return this.ClusterInfo.Status == CLUSTER_STATUS_READY
}

func (this *Cluster) Ready() {
	this.ClusterInfo.Status = CLUSTER_STATUS_READY
}
func (this *Cluster) Join() {
	this.ClusterInfo.Status = CLUSTER_STATUS_JOIN
}

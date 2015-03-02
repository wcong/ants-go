package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"log"
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

// basic cluster infomation
type ClusterInfo struct {
	Name       string
	NodeList   []*NodeInfo
	LocalNode  *NodeInfo
	MasterNode *NodeInfo
}

// recive basic request and record crawled requets
type RequestStatus struct {
	CrawledMap   map[string]int
	CrawlingMap  map[string]map[uint64]*http.Request
	WaitingQuene *crawler.RequestQuene
}

type Cluster struct {
	ClusterInfo   *ClusterInfo
	RequestStatus *RequestStatus
}

func NewCluster(settings *util.Settings, localNode *NodeInfo) *Cluster {
	clusterInfo := &ClusterInfo{settings.Name, make([]*NodeInfo, 0), localNode, localNode}
	requestStatus := &RequestStatus{}
	requestStatus.CrawledMap = make(map[string]int)
	requestStatus.CrawlingMap = make(map[string]map[uint64]*http.Request)
	requestStatus.WaitingQuene = crawler.NewRequestQuene()
	cluster := &Cluster{clusterInfo, requestStatus}
	cluster.ClusterInfo.NodeList = append(cluster.ClusterInfo.NodeList, localNode)
	return cluster
}

// is local node master node
func (this *Cluster) IsMasterNode() bool {
	return this.ClusterInfo.LocalNode.Name == this.ClusterInfo.MasterNode.Name
}

// add a node to cluster node list
func (this *Cluster) AddNode(nodeInfo *NodeInfo) {
	this.ClusterInfo.NodeList = append(this.ClusterInfo.NodeList, nodeInfo)
	if this.ClusterInfo.LocalNode == this.ClusterInfo.MasterNode {
		this.ElectMaster()
	}
}

// choose a new master node
func (this *Cluster) ElectMaster() *NodeInfo {
	this.ClusterInfo.MasterNode = this.ClusterInfo.NodeList[0]
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
		requestMap = make(map[uint64]*http.Request)
		this.RequestStatus.CrawlingMap[request.NodeName] = requestMap
	}
	requestMap[request.UniqueName] = request
}

// a request job is done
// delete it from crawling quene
// add crawled num
func (this *Cluster) Crawled(nodeName string, requestHashName uint64) {
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

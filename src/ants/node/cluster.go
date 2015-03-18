package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"encoding/json"
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
	CrawledMap   map[string]int // node + num
	CrawlingMap  map[string]map[string]*http.Request
	WaitingQuene *crawler.RequestQuene
}

type Cluster struct {
	ClusterInfo   *ClusterInfo
	RequestStatus *RequestStatus
	mutex         *sync.Mutex
	crawlStatus   *crawler.CrawlerStatus
	settings      *util.Settings
}

func NewCluster(settings *util.Settings, localNode *NodeInfo) *Cluster {
	clusterInfo := &ClusterInfo{CLUSTER_STATUS_INIT, settings.Name, make([]*NodeInfo, 0), localNode, nil}
	requestStatus := &RequestStatus{}
	requestStatus.CrawledMap = make(map[string]int)
	requestStatus.CrawlingMap = make(map[string]map[string]*http.Request)
	requestStatus.WaitingQuene = crawler.NewRequestQuene()
	cluster := &Cluster{clusterInfo, requestStatus, new(sync.Mutex), crawler.NewCrawlerStatus(), settings}
	cluster.ClusterInfo.NodeList = append(cluster.ClusterInfo.NodeList, localNode)
	return cluster
}

// get crawl status
func (this *Cluster) CrawlStatus() *crawler.CrawlerStatus {
	return this.crawlStatus
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

//when start a spider ,cluster should record it
func (this *Cluster) StartSpider(spiderName string) {
	this.crawlStatus.StartSpider(spiderName)
}

// pop a request from waiting quene
// add to crawling quenu
func (this *Cluster) PopRequest() *http.Request {
	request := this.RequestStatus.WaitingQuene.Pop()
	return request
}

// record distribute request job
func (this *Cluster) AddToCrawlingQuene(request *http.Request) {
	requestMap, ok := this.RequestStatus.CrawlingMap[request.NodeName]
	if !ok {
		requestMap = make(map[string]*http.Request)
		this.RequestStatus.CrawlingMap[request.NodeName] = requestMap
	}
	this.crawlStatus.Distribute(request.SpiderName)
	requestMap[request.UniqueName] = request
}

// a request job is done
// delete it from crawling quene
// add crawled num
func (this *Cluster) Crawled(nodeName, requestHashName string) {
	requestMap, nodeOk := this.RequestStatus.CrawlingMap[nodeName]
	if !nodeOk {
		log.Println("none node :" + nodeName)
		return
	}
	request, requestOk := requestMap[requestHashName]
	if !requestOk {
		log.Println("none request :" + requestHashName)
		return
	}
	// change RequestStatus
	this.RequestStatus.CrawledMap[nodeName] += 1
	delete(requestMap, requestHashName)
	// change  crawlStatus
	this.crawlStatus.Crawled(request.SpiderName)
	if this.crawlStatus.CanWeStop(request.SpiderName) {
		spiderStatus := this.crawlStatus.CloseSpider(request.SpiderName)
		message, _ := json.Marshal(spiderStatus)
		util.DumpResult(this.settings.LogPath, spiderStatus.Name, string(message))
	}
}

// add a request to quene
func (this *Cluster) AddRequest(request *http.Request) {
	this.RequestStatus.WaitingQuene.Push(request)
	this.crawlStatus.Push(request.SpiderName)
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

// is the spider running
func (this *Cluster) IsSpiderRunning(spiderName string) bool {
	return this.crawlStatus.IsSpiderRunning(spiderName)
}

// get master node
func (this *Cluster) GetMasterNode() *NodeInfo {
	return this.ClusterInfo.MasterNode
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

package support

import (
	"encoding/json"
	. "github.com/wcong/ants-go/ants/cluster"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"github.com/wcong/ants-go/ants/util"
	"log"
	"sync"
)

type DefaultCluster struct {
	clusterInfo   *ClusterInfo
	RequestStatus *RequestStatus
	mutex         *sync.Mutex
	crawlStatus   *crawler.CrawlerStatus
	settings      *util.Settings
}

func NewDefaultCluster(settings *util.Settings) *DefaultCluster {
	clusterInfo := &ClusterInfo{CLUSTER_STATUS_INIT, settings.Name, make([]*node.NodeInfo, 0), nil, nil}
	requestStatus := NewRequestStatus()
	cluster := &DefaultCluster{clusterInfo, requestStatus, new(sync.Mutex), crawler.NewCrawlerStatus(), settings}
	return cluster
}

func (this *DefaultCluster) Init(node *node.NodeInfo) {
	this.clusterInfo.LocalNode = node
	this.clusterInfo.NodeList = append(this.clusterInfo.NodeList, node)
}

func (this *DefaultCluster) GetClusterInfo() *ClusterInfo {
	return this.clusterInfo
}

// get crawl status
func (this *DefaultCluster) CrawlStatus() *crawler.CrawlerStatus {
	return this.crawlStatus
}

func (this *DefaultCluster) DeleteDeadNode(nodeName string) {
	this.mutex.Lock()
	this.removeNode(nodeName)
	this.RequestStatus.DeleteDeadNode(nodeName)
	this.mutex.Unlock()
}

// of course it is not local node
func (this *DefaultCluster) removeNode(nodeName string) {
	deleteIndex := -1
	for index, node := range this.clusterInfo.NodeList {
		if node.Name == nodeName {
			deleteIndex = index
		}
	}
	if deleteIndex >= 0 {
		oldNodeList := this.clusterInfo.NodeList
		this.clusterInfo.NodeList = append(oldNodeList[0:deleteIndex], oldNodeList[deleteIndex+1:len(oldNodeList)]...)
	}
}

// is local node master node
func (this *DefaultCluster) IsMasterNode() bool {
	if this.clusterInfo.MasterNode == nil {
		return false
	}
	return this.clusterInfo.LocalNode.Name == this.clusterInfo.MasterNode.Name
}

// add a node to cluster node list
func (this *DefaultCluster) AddNode(nodeInfo *node.NodeInfo) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, node := range this.clusterInfo.NodeList {
		if node.Name == nodeInfo.Name {
			return
		}
	}
	this.clusterInfo.NodeList = append(this.clusterInfo.NodeList, nodeInfo)
}

func (this *DefaultCluster) GetAllNode() []*node.NodeInfo {
	return this.clusterInfo.NodeList
}

func (this *DefaultCluster) GetRequestStatus() *RequestStatus {
	return this.RequestStatus
}

// make master node by node name
func (this *DefaultCluster) MakeMasterNode(nodeName string) {
	for _, node := range this.clusterInfo.NodeList {
		if node.Name == nodeName {
			this.clusterInfo.MasterNode = node
		}
	}
}

// choose a new master node
func (this *DefaultCluster) ElectMaster() *node.NodeInfo {
	this.clusterInfo.MasterNode = this.clusterInfo.NodeList[0]
	for _, node := range this.clusterInfo.NodeList {
		if this.clusterInfo.MasterNode.Name > node.Name {
			this.clusterInfo.MasterNode = node
		}
	}
	return this.clusterInfo.MasterNode
}

//when start a spider ,cluster should record it
func (this *DefaultCluster) StartSpider(spiderName string) {
	this.crawlStatus.StartSpider(spiderName)
}

// pop a request from waiting quene
// add to crawling quenu
func (this *DefaultCluster) PopRequest() *http.Request {
	request := this.RequestStatus.WaitingQuene.Pop()
	return request
}

// record distribute request job
func (this *DefaultCluster) AddToCrawlingQuene(request *http.Request) {
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
func (this *DefaultCluster) Crawled(scrapyResult *crawler.ScrapeResult) {
	this.RequestStatus.Crawled(scrapyResult)
	// change  crawlStatus
	this.crawlStatus.Crawled(scrapyResult.Request.SpiderName)
}

func (this *DefaultCluster) CanWeStopSpider(spiderName string) bool {
	return this.crawlStatus.CanWeStop(spiderName)
}

func (this *DefaultCluster) StopSpider(spiderName string) {
	spiderStatus := this.crawlStatus.CloseSpider(spiderName)
	message, _ := json.Marshal(spiderStatus)
	log.Println("dump", spiderName, "result")
	util.DumpResult(this.settings.LogPath, spiderStatus.Name, string(message))
}

// add a request to quene
func (this *DefaultCluster) AddRequest(request *http.Request) {
	this.RequestStatus.WaitingQuene.Push(request)
	this.crawlStatus.Push(request.SpiderName)
}

// get master node name
func (this *DefaultCluster) GetMasterName() string {
	return this.clusterInfo.MasterNode.Name
}

// is all loop stop
func (this *DefaultCluster) IsStop() bool {
	return this.RequestStatus.IsStop()
}

func (this *DefaultCluster) HasNode(nodeName string) bool {
	for _, node := range this.clusterInfo.NodeList {
		if node.Name == nodeName {
			return true
		}
	}
	return false
}

// is the spider running
func (this *DefaultCluster) IsSpiderRunning(spiderName string) bool {
	return this.crawlStatus.IsSpiderRunning(spiderName)
}

// get master node
func (this *DefaultCluster) GetMasterNode() *node.NodeInfo {
	return this.clusterInfo.MasterNode
}

// is cluster ready for crawl
func (this *DefaultCluster) IsReady() bool {
	return this.clusterInfo.Status == CLUSTER_STATUS_READY
}

func (this *DefaultCluster) Ready() {
	this.clusterInfo.Status = CLUSTER_STATUS_READY
}

func (this *DefaultCluster) Join() {
	this.clusterInfo.Status = CLUSTER_STATUS_JOIN
}

package support

import (
	"github.com/wcong/ants-go/ants/cluster"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"github.com/wcong/ants-go/ants/util"
	"log"
	"strconv"
)

type DefaultNode struct {
	NodeInfo *node.NodeInfo
	Settings *util.Settings
	Cluster  cluster.Cluster
	Crawler  *crawler.Crawler
}

func NewDefaultNode(settings *util.Settings, resultQuene *crawler.ResultQuene) *DefaultNode {
	ip := util.GetLocalIp()
	name := ip + ":" + strconv.Itoa(settings.TcpPort)
	nodeInfo := &node.NodeInfo{name, ip, settings.TcpPort, settings}
	crawler := crawler.NewCrawler(resultQuene, settings)
	return &DefaultNode{
		NodeInfo: nodeInfo,
		Settings: settings,
		Crawler:  crawler,
	}
}

func (this *DefaultNode) Init(cluster cluster.Cluster) {
	this.Cluster = cluster
}

func (this *DefaultNode) GetNodeInfo() *node.NodeInfo {
	return this.NodeInfo
}

// if spider is running return false
// tell cluster start a spider
// get start requests, push them to cluster request
// try to start the crawler,
func (this *DefaultNode) StartSpider(spiderName string) (bool, string) {
	if this.Cluster.IsSpiderRunning(spiderName) {
		return false, "spider is running"
	}
	log.Println("start spider", spiderName)
	this.Cluster.StartSpider(spiderName)
	this.Crawler.StartSpider(spiderName)
	startRequest := this.Crawler.GetStartRequest(spiderName)
	for _, request := range startRequest {
		this.Cluster.AddRequest(request)
	}
	this.Crawler.Start()
	return true, "started"
}

func (this *DefaultNode) CloseSpider(spiderName string) {
	log.Println("close spider", spiderName)
	this.Crawler.CloseSpider(spiderName)
	this.Cluster.StopSpider(spiderName)
}

// get distribute request
// if node not running ,start it
func (this *DefaultNode) AcceptRequest(request *http.Request) {
	this.Crawler.Downloader.RequestQuene.Push(request)
	this.StartCrawl()
}

// is the node is myself
func (this *DefaultNode) IsMe(nodeName string) bool {
	return this.NodeInfo.Name == nodeName
}

// distribute request to every node
// judge node
// tell cluster where is the request
func (this *DefaultNode) DistributeRequest(request *http.Request) {
	this.Crawler.RequestQuene.Push(request)
	this.AddToCrawlingQuene(request)
}

func (this *DefaultNode) AddToCrawlingQuene(request *http.Request) {
	this.Cluster.AddToCrawlingQuene(request)
}

// report result of request to master
func (this *DefaultNode) ReportToMaster(result *crawler.ScrapeResult) {
	if this.Cluster.IsMasterNode() {
		this.AcceptResult(result)
	}
}

// result of crawl request
// add scraped request to cluster
// tell cluster request is down
func (this *DefaultNode) AcceptResult(scrapyResult *crawler.ScrapeResult) {
	if len(scrapyResult.ScrapedRequests) > 0 {
		for _, request := range scrapyResult.ScrapedRequests {
			if request != nil {
				this.Cluster.AddRequest(request)
			}
		}
	}
	// push request first , avoid spider shut down
	this.Cluster.Crawled(scrapyResult)
}

func (this *DefaultNode) CanWeStopSpider(spiderName string) bool {
	return this.Cluster.CanWeStopSpider(spiderName)
}

// if there is none request left ,return true
func (this *DefaultNode) IsStop() bool {
	return this.Cluster.IsStop()
}

// stop all crawl job
func (this *DefaultNode) StopCrawl() {
	this.Crawler.Stop()
}

// get master name of cluster
func (this *DefaultNode) GetMasterName() string {
	return this.Cluster.GetMasterName()
}

// get master node of cluster
func (this *DefaultNode) GetMasterNode() *node.NodeInfo {
	return this.Cluster.GetMasterNode()
}

// make master node
func (this *DefaultNode) MakeMasterNode(nodeName string) {
	this.Cluster.MakeMasterNode(nodeName)
}

// if this is the master node
func (this *DefaultNode) IsMasterNode() bool {
	return this.Cluster.IsMasterNode()
}

// first of all this is master node
// parse crawler
// remove it from cluster
// unparse crawler
func (this *DefaultNode) DeleteDeadNode(nodeName string) {
	this.PauseCrawl()
	this.Cluster.DeleteDeadNode(nodeName)
	this.UnpauseCrawl()
}

func (this *DefaultNode) Join() {
	this.Cluster.Join()
	this.PauseCrawl()
}

func (this *DefaultNode) Ready() {
	this.Cluster.Ready()
	this.UnpauseCrawl()
}

// start dead loop for all job
func (this *DefaultNode) StartCrawl() {
	go this.Crawler.Start()
}

// pause crawl
func (this *DefaultNode) PauseCrawl() {
	this.Crawler.Pause()
}

// unpause crawl
func (this *DefaultNode) UnpauseCrawl() {
	this.Crawler.UnPause()
}

func (this *DefaultNode) GetSpidersName() []string {
	spiderList := make([]string, 0, len(this.Crawler.SpiderMap))
	for spider := range this.Crawler.SpiderMap {
		spiderList = append(spiderList, spider)
	}
	return spiderList
}

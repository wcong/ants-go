package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"strconv"
)

type NodeInfo struct {
	Name     string
	Ip       string
	Port     int
	Settings *util.Settings
}

type Node struct {
	NodeInfo *NodeInfo
	Settings *util.Settings
	Cluster  *Cluster
	Crawler  *crawler.Crawler
}

func NewNode(settings *util.Settings, resultQuene *crawler.ResultQuene) *Node {
	ip := util.GetLocalIp()
	name := strconv.FormatUint(util.HashString(ip+strconv.Itoa(settings.TcpPort)), 10)
	nodeInfo := &NodeInfo{name, ip, settings.TcpPort, settings}
	crawler := crawler.NewCrawler(resultQuene)
	cluster := NewCluster(settings, nodeInfo)
	return &Node{
		NodeInfo: nodeInfo,
		Settings: settings,
		Cluster:  cluster,
		Crawler:  crawler,
	}
}

// add a node to cluster
// if this is master node,elect a new master node and send it to other
func (this *Node) AddNodeToCluster(nodeInfo *NodeInfo) {
	this.Cluster.AddNode(nodeInfo)
}

// slave node get request of master node info then change the master node
func (this *Node) AddMasterNode(masterNodeInfo *NodeInfo) {
	for _, nodeInfo := range this.Cluster.ClusterInfo.NodeList {
		if nodeInfo.Name == masterNodeInfo.Name {
			this.Cluster.ClusterInfo.MasterNode = nodeInfo
			break
		}
	}
}

// start a spider if not deap loop distribute ,start it
// start a reporter report the crawl result
func (this *Node) StartSpider(spiderName string) (bool, string) {
	if this.Cluster.IsSpiderRunning(spiderName) {
		return false, "spider is running"
	}
	this.Crawler.StartSpider(spiderName)
	this.Cluster.StartSpider(spiderName)
	return true, "started"
}

// get distribute request
// if node not running ,start it
func (this *Node) AcceptRequest(request *http.Request) {
	this.Crawler.Downloader.RequestQuene.Push(request)
	this.StartCrawl()
}

// is the node is myself
func (this *Node) IsMe(nodeName string) bool {
	return this.NodeInfo.Name == nodeName
}

// distribute request to every node
// judge node
// tell cluster where is the request
func (this *Node) DistributeRequest(request *http.Request) {
	this.Crawler.RequestQuene.Push(request)
	this.AddToCrawlingQuene(request)
}

func (this *Node) AddToCrawlingQuene(request *http.Request) {
	this.Cluster.AddToCrawlingQuene(request)
}

// report result of request to master
func (this *Node) ReportToMaster(result *crawler.ScrapeResult) {
	if this.Cluster.IsMasterNode() {
		this.AcceptResult(result)
	}
}

// result of crawl request
// add scraped request to cluster
// tell cluster request is down
// close if cluster has no further request and running request
func (this *Node) AcceptResult(scrapyResult *crawler.ScrapeResult) {
	if len(scrapyResult.ScrapedRequests) > 0 {
		for _, request := range scrapyResult.ScrapedRequests {
			if request != nil {
				this.Cluster.AddRequest(request)
			}
		}
	}
	// push request first , avoid spider shut down
	this.Cluster.Crawled(scrapyResult.Request.NodeName, scrapyResult.Request.UniqueName)
}

// if there is none request left ,return true
func (this *Node) IsStop() bool {
	return this.Cluster.IsStop()
}

// close all node
func (this *Node) GetAllNodeForClose() []*NodeInfo {
	return this.Cluster.ClusterInfo.NodeList
}

// stop all crawl job
func (this *Node) StopCrawl() {
	this.Crawler.StopSpider()
}

// get master name of cluster
func (this *Node) GetMasterName() string {
	return this.Cluster.GetMasterName()
}

// get master node of cluster
func (this *Node) GetMasterNode() *NodeInfo {
	return this.Cluster.GetMasterNode()
}

// make master node
func (this *Node) MakeMasterNode(nodeName string) {
	this.Cluster.MakeMasterNode(nodeName)
}

// if this is the master node
func (this *Node) IsMasterNode() bool {
	return this.Cluster.IsMasterNode()
}

func (this *Node) Join() {
	this.Cluster.Join()
	this.PauseCrawl()
}

func (this *Node) Ready() {
	this.Cluster.Ready()
	this.UnpauseCrawl()
}

// start dead loop for all job
func (this *Node) StartCrawl() {
	go this.Crawler.Start()
}

// pause crawl
func (this *Node) PauseCrawl() {
	this.Crawler.Pause()
}

// unpause crawl
func (this *Node) UnpauseCrawl() {
	this.Crawler.UnPause()
}

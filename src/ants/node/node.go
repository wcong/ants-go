package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"encoding/json"
	"log"
	"strconv"
	"sync"
)

/*
what a node do
*	init a node , a cluster and a crawler
*	load spiders
*	start http service
*	start tcp service
*   accept all http and tcp  request , process it to cluster or crawler
*	*	start a spider
*	*	distribute request
* 	* 	accept request result some request or just scrapy result
*/

type NodeInfo struct {
	Name     string
	Ip       string
	Port     int
	Settings *util.Settings
}

type Node struct {
	NodeInfo    *NodeInfo
	Settings    *util.Settings
	Crawler     *crawler.Crawler
	HttpServer  *http.HttpServer
	Transporter *Transporter
	Cluster     *Cluster
	Distributer *Distributer
	Reporter    *Reporter
}

func NewNode(settings *util.Settings) *Node {
	ip := util.GetLocalIp()
	name := strconv.FormatUint(util.HashString(ip+strconv.Itoa(settings.TcpPort)), 10)
	return &Node{
		NodeInfo: &NodeInfo{
			Name:     name,
			Ip:       ip,
			Port:     settings.TcpPort,
			Settings: settings},
		Settings: settings,
	}
}

// init all base service and container
func (this *Node) Init() {
	this.Reporter = NewReporter(this)
	this.Cluster = NewCluster(this.Settings, this.NodeInfo)
	this.Crawler = crawler.NewCrawler(this.Reporter.ResultQuene)
	this.Crawler.LoadSpiders()
	router := NewRouter(this)
	this.HttpServer = http.NewHttpServer(this.Settings, router)
	transporter := NewTransporter(this.Settings, this)
	this.Transporter = transporter
	this.Distributer = NewDistributer(this.Cluster, this)
}

// start to server
func (this *Node) Start() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go this.HttpServer.Start(wg)
	this.Transporter.Start()
	log.Println("ok,we are ready")
	wg.Wait()
	log.Println("shutting down,goods bye")
}

// add a node to cluster
// if this is master node,elect a new master node and send it to other
func (this *Node) AddNodeToCluster(nodeInfo *NodeInfo) {
	this.Cluster.AddNode(nodeInfo)
	//if this.Cluster.ClusterInfo.LocalNode == this.Cluster.ClusterInfo.MasterNode {
	//	masterNode := this.Cluster.ElectMaster()
	//	jsonMessage := RequestMessage{
	//		Type:     HANDLER_SEND_MASTER_REQUEST,
	//		NodeInfo: masterNode,
	//	}
	//	json, _ := json.Marshal(jsonMessage)
	//	message := string(json)
	//	for _, nodeInfo := range this.Cluster.ClusterInfo.NodeList {
	//		if nodeInfo.Name == this.NodeInfo.Name {
	//			continue
	//		}
	//		this.Transporter.SendMessageToNode(nodeInfo.Name, message)
	//	}
	//}
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
func (this *Node) StartSpider(spiderName string) *crawler.StartSpiderResult {
	result := this.Crawler.StartSpider(spiderName)
	if result.Success && this.Distributer.IsStop() {
		go this.Distributer.Start()
	}
	if result.Success && this.Reporter.IsStop() {
		go this.Reporter.Start()
	}
	return result
}

// start dead loop for all job
func (this *Node) StartCrawl() {
	go this.Distributer.Start()
	go this.Reporter.Start()
	go this.Crawler.Start()
}

// get distribute request
// if node not running ,start it
func (this *Node) AcceptRequest(request *http.Request) {
	this.Crawler.Downloader.RequestQuene.Push(request)
	this.StartCrawl()
}

// distribute request to every node
// judge node
// tell cluster where is the request
func (this *Node) DistributeRequest(request *http.Request) {
	if request.NodeName == this.NodeInfo.Name {
		this.Crawler.RequestQuene.Push(request)
		this.Cluster.AddToCrawlingQuene(request)
	} else {
		requestSlice := make([]*http.Request, 1)
		requestSlice[0] = request
		jsonMessage := RequestMessage{
			Type:    HANDLER_SEND_REQUEST,
			Request: request,
		}
		message, err := json.Marshal(jsonMessage)
		if err != nil {
			log.Println(err)
		} else {
			this.Transporter.SendMessageToNode(request.NodeName, string(message))
			this.Cluster.AddToCrawlingQuene(request)
		}
	}
}

// report result of request to master
func (this *Node) ReportToMaster(result *crawler.ScrapeResult) {
	requestMessage := &RequestMessage{
		Type:            HANDLER_SEND_REQUEST_RESULT,
		Request:         result.Request,
		CrawlResult:     result.CrawlResult,
		ScrapedRequests: result.ScrapedRequests,
		NodeInfo:        this.NodeInfo,
	}
	if this.Cluster.IsMasterNode() {
		this.AcceptResult(requestMessage)
		return
	}
	message, err := json.Marshal(requestMessage)
	if err != nil {
		log.Panic(err)
		return
	}
	this.Transporter.SendMessageToNode(this.Cluster.GetMasterName(), string(message))
}

// result of crawl request
// tell cluster request is down
// add scraped request to cluster
// close if cluster has no further request and running request
func (this *Node) AcceptResult(responseMessage *RequestMessage) {
	this.Cluster.Crawled(responseMessage.Request.NodeName, responseMessage.Request.UniqueName)
	if len(responseMessage.ScrapedRequests) > 0 {
		for _, request := range responseMessage.ScrapedRequests {
			this.Cluster.AddRequest(request)
		}
	}
	if this.Cluster.IsStop() {
		this.CloseAllNode()
	}
}

// close all node
func (this *Node) CloseAllNode() {
	requestMessage := &RequestMessage{}
	requestMessage.Type = HANDLER_STOP_NODE
	json, _ := json.Marshal(requestMessage)
	message := string(json)
	for _, nodeInfo := range this.Cluster.ClusterInfo.NodeList {
		if nodeInfo.Name == this.NodeInfo.Name {
			this.StopCrawl()
			continue
		}
		this.Transporter.SendMessageToNode(nodeInfo.Name, message)
	}
}

// stop all crawl job
func (this *Node) StopCrawl() {
	this.Crawler.StopSpider()
	this.Distributer.Stop()
	this.Reporter.Stop()
}

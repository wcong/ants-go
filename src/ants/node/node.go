package node

import (
	"ants/crawler"
	"ants/http"
	"ants/util"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"
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
	this.Cluster = NewCluster(this.Settings, this.NodeInfo)
	this.Crawler = crawler.NewCrawler()
	this.Crawler.LoadSpiders()
	router := NewRouter(this)
	this.HttpServer = http.NewHttpServer(this.Settings, router)
	transporter := NewTransporter(this.Settings, this)
	this.Transporter = transporter
	this.Distributer = NewDistributer(this.Cluster)
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
	if this.Cluster.ClusterInfo.LocalNode == this.Cluster.ClusterInfo.MasterNode {
		masterNode := this.Cluster.ElectMaster()
		jsonMessage := JSONMessage{
			Type:     HANDLER_SEND_MASTER_REQUEST,
			NodeInfo: *masterNode,
		}
		json, _ := json.Marshal(jsonMessage)
		message := string(json)
		for _, nodeInfo := range this.Cluster.ClusterInfo.NodeList {
			if nodeInfo.Name == this.NodeInfo.Name {
				continue
			}
			this.Transporter.SendMessageToNode(nodeInfo.Name, message)
		}
	}
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
func (this *Node) StartSpider(spiderName string) *crawler.StartSpiderResult {
	result := this.Crawler.StartSpider(spiderName)
	if result.Success && this.Distributer.IsStop() {
		this.Distributer.Run()
		go this.DistributeRequest()
	}
	return result
}

// dead loop cluster pop request
// distribute request to every node
// judge node
// tell cluster where is the request
func (this *Node) DistributeRequest() {
	for {
		if this.Distributer.IsParse() {
			time.Sleep(1 * time.Second)
			continue
		}
		request := this.Cluster.PopRequest()
		if request == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		nodeName := this.Distributer.Distribute(request)
		if nodeName == this.NodeInfo.Name {
			this.Crawler.RequestQuene.Push(request)
			this.Cluster.AddToCrawlingQuene(request)
		} else {
			jsonMessage := JSONMessage{
				Type: HANDLER_SEND_REQUEST,
			}
			message, err := json.Marshal(jsonMessage)
			if err != nil {
				log.Println(err)
			} else {
				this.Transporter.SendMessageToNode(nodeName, string(message))
				this.Cluster.AddToCrawlingQuene(request)
			}
		}
	}
}

// result of crawl request
// if
func (this *Node) AcceptResult(jsonMessage *JSONMessage) {

}

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
func (this *Node) Init() {
	this.Crawler = crawler.NewCrawler()
	this.Crawler.LoadSpiders()
	router := NewRouter(this)
	this.HttpServer = http.NewHttpServer(this.Settings, router)
	this.Cluster = NewCluster(this.Settings, this.NodeInfo)
	transporter := NewTransporter(this.Settings, this)
	this.Transporter = transporter
}

func (this *Node) Start() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go this.HttpServer.Start(wg)
	this.Transporter.Start()
	log.Println("ok,we are ready")
	wg.Wait()
	log.Println("shutting down")
}
func (this *Node) AddNodeToCluster(nodeInfo *NodeInfo) {
	this.Cluster.AddNode(nodeInfo)
	if this.Cluster.LocalNode == this.Cluster.MasterNode {
		masterNode := this.Cluster.ElectMaster()
		jsonMessage := JSONMessage{
			Type:     HANDLER_SEND_MASTER_REQUEST,
			NodeInfo: *masterNode,
		}
		json, _ := json.Marshal(jsonMessage)
		message := string(json)
		for _, nodeInfo := range this.Cluster.NodeList {
			if nodeInfo.Name == this.NodeInfo.Name {
				continue
			}
			this.Transporter.SendMessageToNode(nodeInfo, message)
		}
	}
}

// slave node get request of master node info then change the master node
func (this *Node) AddMasterNode(masterNodeInfo *NodeInfo) {
	for _, nodeInfo := range this.Cluster.NodeList {
		if nodeInfo.Name == masterNodeInfo.Name {
			this.Cluster.MasterNode = nodeInfo
			break
		}
	}
}

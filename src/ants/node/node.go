package node

import (
	"ants/conf"
	"ants/crawler"
	"ants/http"
	"ants/util"
	"log"
	"strconv"
	"sync"
)

type NodeInfo struct {
	Name     string
	Ip       string
	Port     int
	Settings *conf.Settings
}

type Node struct {
	NodeInfo    *NodeInfo
	Settings    *conf.Settings
	Crawler     *crawler.Crawler
	HttpServer  *http.HttpServer
	Transporter *Transporter
	Cluster     *Cluster
}

func NewNode(settings *conf.Settings) *Node {
	ip := util.GetLocalIp()
	name := strconv.FormatUint(util.HashString(ip+strconv.Itoa(settings.TcpPort)), 10)
	return &Node{
		NodeInfo: &NodeInfo{
			Name:     name,
			Ip:       ip,
			Port:     settings.TcpPort,
			Settings: settings},
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
	log.Println("ok,we are ready")
	wg.Wait()
}
func (this *Node) AddNodeToCluster(nodeInfo *NodeInfo) {

}

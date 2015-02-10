package node

import (
	"ants/conf"
	"ants/crawler"
	"ants/http"
	"ants/transport"
	"ants/util"
	"strconv"
)

type Node struct {
	Name       string
	Ip         string
	Port       int
	Settings   *conf.Settings
	Crawler    *crawler.Crawler
	HttpServer *http.HttpServer
	TcpServer  *transport.TcpManager
}

func NewNode(settings *conf.Settings) *Node {
	ip := util.GetLocalIp()
	name := strconv.FormatUint(util.HashString(ip+strconv.Itoa(settings.TcpPort)), 10)
	return &Node{
		Name:     name,
		Ip:       ip,
		Port:     settings.TcpPort,
		Settings: settings,
	}
}
func (this *Node) Init() {
	initChan := make(chan int)
	this.Crawler = crawler.Crawler{}
	this.Crawler.LoadSpiders()
	http.InitServer(initChan)
}

func (this *Node) Start() {
}

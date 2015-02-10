package node

import (
	"ants/conf"
	"ants/crawler"
	"ants/http"
	"ants/transport"
	"ants/util"
	"log"
	"strconv"
	"sync"
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
	this.Crawler = &crawler.Crawler{}
	this.Crawler.LoadSpiders()
	router := NewRouter(this)
	this.HttpServer = http.NewHttpServer(this.Settings, router)
}

func (this *Node) Start() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go this.HttpServer.Start(wg)
	log.Println("ok,we are ready")
	wg.Wait()
}

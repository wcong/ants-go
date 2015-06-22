package node

import (
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/util"
)

type NodeInfo struct {
	Name     string
	Ip       string
	Port     int
	Settings *util.Settings
}

type Node interface {
	GetNodeInfo() *NodeInfo
	StartSpider(spiderName string) (bool, string)
	CloseSpider(spiderName string)
	AcceptRequest(request *http.Request)
	IsMe(nodeName string) bool
	DistributeRequest(request *http.Request)
	AddToCrawlingQuene(request *http.Request)
	ReportToMaster(result *crawler.ScrapeResult)
	AcceptResult(scrapyResult *crawler.ScrapeResult)
	CanWeStopSpider(spiderName string) bool
	IsStop() bool
	StopCrawl()
	MakeMasterNode(nodeName string)
	IsMasterNode() bool
	Join()
	Ready()
	StartCrawl()
	PauseCrawl()
	UnpauseCrawl()
	GetSpidersName() []string
}

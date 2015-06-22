package cluster

import (
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
)

// cluster status
// *		init:where every thing have init
// *		join:try to connect to other node ,if not make itself master,else ,get other master
// *		election(option):when circle is builded a start to elect a master
// * 		ready:ready to start crawl
const (
	CLUSTER_STATUS_INIT = iota
	CLUSTER_STATUS_JOIN
	CLUSTER_STATUS_ELECTION
	CLUSTER_STATUS_READY
)

// basic cluster infomation
type ClusterInfo struct {
	Status     int
	Name       string
	NodeList   []*node.NodeInfo
	LocalNode  *node.NodeInfo
	MasterNode *node.NodeInfo
}

type Cluster interface {
	CrawlStatus() *crawler.CrawlerStatus
	DeleteDeadNode(nodeName string)
	IsMasterNode() bool
	AddNode(nodeInfo *node.NodeInfo)
	MakeMasterNode(nodeName string)
	ElectMaster() *node.NodeInfo
	IsReady() bool
	IsSpiderRunning(spiderName string) bool
	StartSpider(spiderName string)
	AddRequest(request *http.Request)
	StopSpider(spiderName string)
	AddToCrawlingQuene(request *http.Request)
	Crawled(scrapyResult *crawler.ScrapeResult)
	CanWeStopSpider(spiderName string) bool
	IsStop() bool
	PopRequest() *http.Request
	Ready()
	Join()
	GetMasterNode() *node.NodeInfo
	GetMasterName() string
	GetClusterInfo() *ClusterInfo
	GetRequestStatus() *RequestStatus
	GetAllNode() []*node.NodeInfo
}

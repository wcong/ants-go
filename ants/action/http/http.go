package http

import (
	"encoding/json"
	"github.com/wcong/ants-go/ants/action"
	"github.com/wcong/ants-go/ants/cluster"
	Http "github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	"log"
	"net/http"
	"time"
)

type Router struct {
	node        node.Node
	cluster     cluster.Cluster
	mux         map[string]func(http.ResponseWriter, *http.Request)
	reporter    action.Watcher
	distributer action.Watcher
	rpcClient   action.RpcClientAnts
}

func NewRouter(node node.Node, cluster cluster.Cluster, reporter, distributer action.Watcher, rpcClient action.RpcClientAnts) *Router {
	mux := make(map[string]func(http.ResponseWriter, *http.Request))
	router := &Router{
		node:        node,
		cluster:     cluster,
		mux:         mux,
		reporter:    reporter,
		distributer: distributer,
		rpcClient:   rpcClient,
	}
	mux["/"] = router.Welcome
	mux["/cluster"] = router.Cluster
	mux["/spiders"] = router.Spiders
	mux["/crawl"] = router.Crawl
	mux["/crawl/cluster"] = router.CrawlCluster
	mux["/crawl/node"] = router.CrawlNode
	return router
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Println("get request:" + url)
	if !this.cluster.IsReady() {
		w.Write([]byte("sorry,cluster not ready,please wait"))
		return
	}
	path := r.URL.Path
	if h, ok := this.mux[path]; ok {
		w.Header().Set(Http.CONTENT_TYPE, Http.JSON_CONTENT_TYPE)
		h(w, r)
		return
	}
	this.Welcome(w, r)
}

func (this *Router) Welcome(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02 15:04:05")
	welcome := WelcomeInfo{
		"for crawl",
		"do not panic",
		now,
	}
	encoder, err := json.Marshal(welcome)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}
func (this *Router) Spiders(w http.ResponseWriter, r *http.Request) {
	spiderList := this.node.GetSpidersName()
	encoder, err := json.Marshal(spiderList)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

// try to start spider
// if ok
// *		tell other node start spider
// *		start reporter and distribute in this node
func (this *Router) Crawl(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	spiderName := r.Form["spider"][0]
	now := time.Now().Format("2006-01-02 15:04:05")
	startResult := &StartSpiderResult{}
	startResult.Time = now
	startResult.Spider = spiderName
	startResult.MasterNode = this.cluster.GetMasterNode()
	if !this.node.IsMasterNode() {
		startResult.Success = false
		startResult.Message = "not master,see master node"
	} else {
		result, message := this.node.StartSpider(spiderName)
		if result {
			log.Println("start spider:", spiderName)
			for _, nodeInfo := range this.cluster.GetAllNode() {
				if !this.node.IsMe(nodeInfo.Name) {
					this.rpcClient.StartSpider(nodeInfo.Name, spiderName)
				}
			}
			go this.reporter.Start()
			go this.distributer.Start()
		}
		startResult.Success = result
		startResult.Message = message
	}
	encoder, err := json.Marshal(startResult)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) Cluster(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.cluster.GetClusterInfo())
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlCluster(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.cluster.GetRequestStatus())
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlNode(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.cluster.CrawlStatus())
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

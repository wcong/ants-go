package http

import (
	"ants/action"
	Http "ants/http"
	Node "ants/node"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Router struct {
	node        *Node.Node
	mux         map[string]func(http.ResponseWriter, *http.Request)
	reporter    action.Watcher
	distributer action.Watcher
}

func NewRouter(node *Node.Node, reporter, distributer action.Watcher) *Router {
	mux := make(map[string]func(http.ResponseWriter, *http.Request))
	router := &Router{
		node:        node,
		mux:         mux,
		reporter:    reporter,
		distributer: distributer,
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
	if !this.node.Cluster.IsReady() {
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
	spiderList := make([]string, 0, len(this.node.Crawler.SpiderMap))
	for spider := range this.node.Crawler.SpiderMap {
		spiderList = append(spiderList, spider)
	}
	encoder, err := json.Marshal(spiderList)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) Crawl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	spiderName := r.Form["spider"][0]
	now := time.Now().Format("2006-01-02 15:04:05")
	startResult := &StartSpiderResult{}
	result, message := this.node.StartSpider(spiderName)
	startResult.Time = now
	startResult.Success = result
	startResult.Spider = spiderName
	startResult.Message = message
	encoder, err := json.Marshal(startResult)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) Cluster(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.node.Cluster)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlCluster(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.node.Cluster.RequestStatus)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlNode(w http.ResponseWriter, r *http.Request) {
	encoder, err := json.Marshal(this.node.Cluster.CrawlStatus())
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

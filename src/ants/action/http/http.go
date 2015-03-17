package action

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
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = router.Welcome
	mux["/cluster"] = router.Cluster
	mux["/spiders"] = router.Spiders
	mux["/crawl"] = router.Crawl
	mux["/crawl/cluster"] = router.CrawlCluster
	mux["/crawl/node"] = router.CrawlNode
	router := &Router{
		node:        node,
		mux:         mux,
		reporter:    reporter,
		distributer: distributer,
	}
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
	if h, ok := this.Mux[path]; ok {
		h(w, r)
		return
	}
	this.Welcome(w, r)
}

type WelcomeInfo struct {
	Message  string
	Greeting string
	Time     string
}

func (this *Router) Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(Http.CONTENT_TYPE, Http.JSON_CONTENT_TYPE)
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
	w.Header().Set("Content-Type", "application/json")
	spiderList := make([]string, 0, len(this.Node.Crawler.SpiderMap))
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
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	spiderName := r.Form["spider"][0]
	now := time.Now().Format("2006-01-02 15:04:05")
	result := this.node.StartSpider(spiderName)
	if result.Success {
		if this.distributer.IsStop() {
			go this.Distributer.Start()
		}
		if this.reporter.IsStop() {
			go this.Reporter.Start()
		}
	}
	result.Time = now
	encoder, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) Cluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder, err := json.Marshal(this.Node.Cluster)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder, err := json.Marshal(this.Node.Cluster.RequestStatus)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

func (this *Router) CrawlNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder, err := json.Marshal(this.Node.Cluster.crawlStatus)
	if err != nil {
		log.Println(err)
	}
	w.Write(encoder)
}

package node

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Router struct {
	Node *Node
	Mux  map[string]func(http.ResponseWriter, *http.Request)
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URl.
		log.Println("get request:" + url)
	if h, ok := this.Mux[url]; ok {
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
	w.Header().Set("Content-Type", "application/json")
	now := time.Now().Format("2006-01-02 15:04:05")
	welcome := WelcomeInfo{
		"for crawl",
		"do not panic",
		now,
	}
	encoder, err := json.Marshal(welcome)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}
func (this *Router) Spiders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spiderList := make([]string, 0, len(this.Node.Crawler.SpiderMap))
	for spider := range this.Node.Crawler.SpiderMap {
		spiderList = append(spiderList, spider)
	}
	encoder, err := json.Marshal(spiderList)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}

type CrawlInfo struct {
	SpiderName string
	Time       string
}

func (this *Router) Crawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spiderName := r.Form["spider"][0]
	now := time.Now().Format("2006-01-02 15:04:05")
	crawlInfo := CrawlInfo{
		spiderName,
		now,
	}
	encoder, err := json.Marshal(crawlInfo)
	if err != nil {
		log.Fatal(err)
	}
	go this.Node.Crawler.StartSpider(spiderName)
	w.Write(encoder)
}

func NewRouter(node *Node) *Router {
	router := &Router{}
	router.Node = node
	router.Mux = make(map[string]func(http.ResponseWriter, *http.Request))
	router.Mux["/"] = router.Welcome
	router.Mux["/spiders"] = router.Spiders
	router.Mux["/crawl"] = router.Crawl
	return router
}

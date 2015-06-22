package cluster

import (
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"log"
	"sync"
)

var mutex sync.Mutex

// receive basic request and record crawled requets
type RequestStatus struct {
	CrawledMap   map[string]int // node + num
	CrawlingMap  map[string]map[string]*http.Request
	WaitingQuene *crawler.RequestQuene
}

func NewRequestStatus() *RequestStatus {
	requestStatus := &RequestStatus{}
	requestStatus.CrawledMap = make(map[string]int)
	requestStatus.CrawlingMap = make(map[string]map[string]*http.Request)
	requestStatus.WaitingQuene = crawler.NewRequestQuene()
	return requestStatus
}

// delete in  CrawlingMap
// add for CrawledMap
func (this *RequestStatus) Crawled(scrapyResult *crawler.ScrapeResult) {
	requestMap, nodeOk := this.CrawlingMap[scrapyResult.Request.NodeName]
	if !nodeOk {
		log.Println("none node :" + scrapyResult.Request.NodeName)
		return
	}
	_, requestOk := requestMap[scrapyResult.Request.UniqueName]
	if !requestOk {
		log.Println("none request :" + scrapyResult.Request.UniqueName)
		return
	}
	// change RequestStatus
	mutex.Lock()
	this.CrawledMap[scrapyResult.Request.NodeName] += 1
	delete(requestMap, scrapyResult.Request.UniqueName)
	mutex.Unlock()
}

// remove request from crawlingmap for dead node
// add those requests to waiting quenu
func (this *RequestStatus) DeleteDeadNode(nodeName string) {
	crawlingMap := this.CrawlingMap[nodeName]
	for _, request := range crawlingMap {
		this.WaitingQuene.Push(request)
	}
	delete(this.CrawlingMap, nodeName)
}

// is all loop stop
func (this *RequestStatus) IsStop() bool {
	if !this.WaitingQuene.IsEmpty() {
		return false
	}
	for _, requestMap := range this.CrawlingMap {
		if len(requestMap) > 0 {
			return false
		}
	}
	return true
}

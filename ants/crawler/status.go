package crawler

import (
	"time"
)

type SpiderStatus struct {
	Name      string
	Crawled   int
	Running   int
	Waiting   int
	StartTime time.Time
	EndTime   time.Time
}

func NewSpiderStatus(name string) *SpiderStatus {
	startTime := time.Now()
	spiderStatus := &SpiderStatus{
		Name:      name,
		StartTime: startTime,
		Crawled:   0,
		Running:   0,
		Waiting:   0,
	}
	return spiderStatus
}

// status of crawler
// crawled spiders and running spiders
type CrawlerStatus struct {
	CrawledSpider []*SpiderStatus
	RunningSpider map[string]*SpiderStatus
}

func NewCrawlerStatus() *CrawlerStatus {
	crawledSpider := make([]*SpiderStatus, 0)
	runningSpider := make(map[string]*SpiderStatus)
	return &CrawlerStatus{crawledSpider, runningSpider}
}

func (this *CrawlerStatus) IsSpiderRunning(spiderName string) bool {
	_, ok := this.RunningSpider[spiderName]
	return ok
}

// add a spider to running map
func (this *CrawlerStatus) StartSpider(spiderName string) {
	_, ok := this.RunningSpider[spiderName]
	if !ok {
		this.RunningSpider[spiderName] = NewSpiderStatus(spiderName)
	}
}

// add a request to wait in spiderName
func (this *CrawlerStatus) Push(spiderName string) {
	this.RunningSpider[spiderName].Waiting += 1
}

// if cluster distribute a request,
// waiting -1 runing +1
func (this *CrawlerStatus) Distribute(spiderName string) {
	spiderStatus, _ := this.RunningSpider[spiderName]
	spiderStatus.Waiting -= 1
	spiderStatus.Running += 1
}

// get crawl result
// runing -1 crawled +1
func (this *CrawlerStatus) Crawled(spiderName string) {
	spiderStatus, _ := this.RunningSpider[spiderName]
	spiderStatus.Running -= 1
	spiderStatus.Crawled += 1
}

// judge a is a spider can stop
func (this *CrawlerStatus) CanWeStop(spiderName string) bool {
	spiderStatus, _ := this.RunningSpider[spiderName]
	leftNum := spiderStatus.Running + spiderStatus.Waiting
	return leftNum <= 0
}

// no more request for spider ,close it
// remove from runningSpider
// add to crawledSpider
func (this *CrawlerStatus) CloseSpider(spiderName string) *SpiderStatus {
	spiderStatus, _ := this.RunningSpider[spiderName]
	spiderStatus.EndTime = time.Now()
	this.CrawledSpider = append(this.CrawledSpider, spiderStatus)
	delete(this.RunningSpider, spiderName)
	return spiderStatus
}

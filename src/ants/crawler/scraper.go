package crawler

import (
	"ants/http"
	"ants/spiders"
	"log"
	"time"
)

const (
	SCRAPY_STATUS_STOP = iota
	SCRAPY_STATUS_RUNING
	SCRAPY_STATUS_PAUSE
)

type ScrapeResult struct {
	Request         *http.Request
	CrawlResult     string // if success just empty string,or error reason
	ScrapedRequests []*http.Request
}

type Scraper struct {
	Status        int
	ResultQuene   *ResultQuene
	ResponseQuene *ResponseQuene
	SpiderMap     map[string]*spiders.Spider
}

//
func NewScraper(resultQuene *ResultQuene, responseQuene *ResponseQuene, spiderMap map[string]*spiders.Spider) *Scraper {
	return &Scraper{SCRAPY_STATUS_STOP, resultQuene, responseQuene, spiderMap}
}

func (this *Scraper) Start() {
	if this.Status == SCRAPY_STATUS_RUNING {
		return
	}
	log.Println("start scraper")
	this.Status = SCRAPY_STATUS_RUNING
	this.Scrapy()
}

func (this *Scraper) Stop() {
	this.Status = SCRAPY_STATUS_STOP
}
func (this *Scraper) Pause() {
	this.Status = SCRAPY_STATUS_PAUSE
}
func (this *Scraper) UnPause() {
	this.Status = SCRAPY_STATUS_RUNING
}

// dead loop for scrapy
// pop a response
// scrapy it
// if scrapy some request, push it to quene
func (this *Scraper) Scrapy() {
	for {
		if this.Status == SCRAPY_STATUS_PAUSE {
			time.Sleep(1 * time.Second)
			continue
		}
		if this.Status != SCRAPY_STATUS_RUNING {
			break
		}
		response := this.ResponseQuene.Pop()
		if response == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println("scrapy :" + response.GoResponse.Request.URL.String())
		requestList, err := this.SpiderMap[response.SpiderName].ParseMap[response.ParserName](response)
		scrapeResult := &ScrapeResult{}
		scrapeResult.Request = response.Request
		if err != nil {
			log.Println(err)
			scrapeResult.CrawlResult = err.Error()
		}
		if requestList != nil {
			scrapeResult.ScrapedRequests = requestList
		}
		this.ResultQuene.Push(scrapeResult)
	}
}

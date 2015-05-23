package crawler

import (
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"log"
	"strconv"
	"time"
)

const (
	SCRAPY_STATUS_STOP = iota
	SCRAPY_STATUS_STOPED
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
	return &Scraper{SCRAPY_STATUS_STOPED, resultQuene, responseQuene, spiderMap}
}

func (this *Scraper) Start() {
	if this.Status == SCRAPY_STATUS_RUNING {
		return
	}
	for {
		if this.Status == SCRAPY_STATUS_STOPED {
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("start scraper")
	this.Status = SCRAPY_STATUS_RUNING
	this.Scrapy()
}

func (this *Scraper) Stop() {
	this.Status = SCRAPY_STATUS_STOP
}
func (this *Scraper) Pause() {
	if this.Status == SCRAPY_STATUS_RUNING {
		this.Status = SCRAPY_STATUS_PAUSE
	}
}
func (this *Scraper) UnPause() {
	if this.Status == SCRAPY_STATUS_PAUSE {
		this.Status = SCRAPY_STATUS_RUNING
	}
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
			this.Status = SCRAPY_STATUS_STOPED
			break
		}
		response := this.ResponseQuene.Pop()
		if response == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		go this.scrapyAndPush(response)
	}
}

// scrapy and push in go routine
func (this *Scraper) scrapyAndPush(response *http.Response) {
	log.Println(response.SpiderName, ":start to scrapy:", response.Request.GoRequest.URL.String())
	defer func() {
		err := recover()
		if err != nil {
			log.Println(response.SpiderName, " ", err)
			scrapeResult := &ScrapeResult{}
			scrapeResult.Request = response.Request
			scrapeResult.CrawlResult = err.(error).Error()
			this.ResultQuene.Push(scrapeResult)
		}
	}()
	requestList, err := this.SpiderMap[response.SpiderName].ParseMap[response.ParserName](response)
	scrapeResult := &ScrapeResult{}
	scrapeResult.Request = response.Request
	if err != nil {
		log.Println(err)
		scrapeResult.CrawlResult = err.Error()
	}
	if requestList != nil {
		for _, request := range requestList {
			request.Depth = response.Request.Depth + 1
		}
		scrapeResult.ScrapedRequests = requestList
		log.Println(response.SpiderName, ":scrapyed:", strconv.Itoa(len(requestList)), "requests from:", response.GoResponse.Request.URL.String())
	}
	this.ResultQuene.Push(scrapeResult)
}

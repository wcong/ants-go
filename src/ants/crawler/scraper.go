package crawler

import (
	"ants/util"
	"log"
	"time"
)

const (
	SCRAPY_STATUS_STOP = iota
	SCRAPY_STATUS_RUNING
	SCRAPY_STATUS_PAUSE
)

type Scraper struct {
	Status        int
	RequestQuene  *RequestQuene
	ResponseQuene *ResponseQuene
	SpiderMap     map[string]*util.Spider
}

func NewScraper(requestQuene *RequestQuene, responseQuene *ResponseQuene, spiderMap map[string]*util.Spider) *Scraper {
	return &Scraper{SCRAPY_STATUS_STOP, requestQuene, responseQuene, spiderMap}
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
		if err != nil {
			log.Fatal(err)
		}
		if requestList != nil {
			for _, request := range requestList {
				this.RequestQuene.Push(request)
			}
		}
	}
}

package crawler

import (
	base_spider "ants/spiders"
	"log"
	"spiders"
)

// crawler
type Crawler struct {
	SpiderMap     map[string]*base_spider.Spider //contains all spiders
	RequestQuene  *RequestQuene                  //all waiting request
	ResponseQuene *ResponseQuene                 //all waiting response for scrape
	Downloader    *Downloader                    //download tools
	Scraper       *Scraper                       //scrape tools
}

// resultQuene is for reporter,make sure it is the same ppointer
func NewCrawler(resultQuene *ResultQuene) *Crawler {
	requestQuene := NewRequestQuene()
	responseQuene := NewResponseQuene()
	spiderMap := spiders.LoadAllSpiders()
	downloader := NewDownloader(requestQuene, responseQuene)
	scraper := NewScraper(resultQuene, responseQuene, spiderMap)
	crawler := &Crawler{spiderMap, requestQuene, responseQuene, downloader, scraper}
	return crawler
}

// start a spider
func (this *Crawler) StartSpider(spiderName string) {
	log.Println("start to crawl spider " + spiderName)
	spider := this.SpiderMap[spiderName]
	startRequests := spider.MakeStartRequests()
	for _, request := range startRequests {
		this.RequestQuene.Push(request)
	}
	this.Start()
}

func (this *Crawler) Start() {
	go this.Downloader.Start()
	go this.Scraper.Start()
}

func (this *Crawler) Pause() {
	this.Downloader.Pause()
	this.Scraper.Pause()
}

func (this *Crawler) UnPause() {
	this.Downloader.UnPause()
	this.Scraper.UnPause()
}

func (this *Crawler) StopSpider() {
	this.Downloader.Stop()
	this.Scraper.Stop()
}

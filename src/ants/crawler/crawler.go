package crawler

import (
	base_spider "ants/spiders"
	"log"
	"spiders"
)

// crawler
// *		contains all spiders
// *		record crawl status
// *		all waiting request
// *		all waiting response for scrape
// *		download tools
// *		scrape tools
type Crawler struct {
	SpiderMap     map[string]*base_spider.Spider
	RequestQuene  *RequestQuene
	ResponseQuene *ResponseQuene
	Downloader    *Downloader
	Scraper       *Scraper
}

func NewCrawler(resultQuene *ResultQuene) *Crawler {
	requestQuene := NewRequestQuene()
	responseQuene := NewResponseQuene()
	spiderMap := spiders.LoadAllSpiders()
	downloader := NewDownloader(requestQuene, responseQuene)
	scraper := NewScraper(resultQuene, responseQuene, spiderMap)
	crawler := Crawler{spiderMap, requestQuene, responseQuene, downloader, scraper}
	return &crawler
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

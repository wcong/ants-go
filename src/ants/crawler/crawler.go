package crawler

import (
	"ants/http"
	base_spider "ants/spiders"
	"log"
	"spiders"
)

type Crawler struct {
	SpiderMap     map[string]*base_spider.Spider
	Status        CrawlerStatus
	RequestQuene  *RequestQuene
	ResponseQuene *ResponseQuene
	Downloader    *Downloader
	Scraper       *Scraper
}

type StartSpiderResult struct {
	Success bool
	Message string
	Spider  string
	Time    string
}

func NewCrawler() *Crawler {
	requestQuene := NewRequestQuene()
	responseQuene := NewResponseQuene()
	spiderMap := make(map[string]*base_spider.Spider)
	downloader := NewDownloader(requestQuene, responseQuene)
	scraper := NewScraper(requestQuene, responseQuene, spiderMap)
	crawler := Crawler{spiderMap, CrawlerStatus{}, requestQuene, responseQuene, downloader, scraper}
	return &crawler
}

func (this *Crawler) LoadSpiders() {
	deadLoopTest := spiders.MakeDeadLoopSpider()
	this.SpiderMap[deadLoopTest.Name] = deadLoopTest
}
func (this *Crawler) StartSpider(spiderName string) *StartSpiderResult {
	log.Println("start to crawl spider " + spiderName)
	spider := this.SpiderMap[spiderName]
	result := &StartSpiderResult{}
	if spider.Status == base_spider.SPIDERS_STATUS_RUNNING {
		result.Success = false
		result.Message = "spider already runing"
		result.Spider = spider.Name
		return result
	}
	spider.Status = base_spider.SPIDERS_STATUS_RUNNING
	for _, url := range spider.StartUrls {
		request, err := http.NewRequest("GET", url, nil, spider.Name, base_spider.BASE_PARSE_NAME)
		if err != nil {
			log.Fatal(err)
			continue
		}
		this.RequestQuene.Push(request)
	}
	this.RunSpider()
	result.Success = true
	result.Message = "start spider"
	result.Spider = spider.Name
	return result
}
func (this *Crawler) RunSpider() {
	go this.Downloader.Start()
	go this.Scraper.Start()
}
func (this *Crawler) ParseSpider() {
	this.Downloader.Pause()
	this.Scraper.Pause()
}
func (this *Crawler) UnParseSpider() {
	this.Downloader.UnPause()
	this.Scraper.UnPause()
}

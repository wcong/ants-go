package crawler

import (
	"ants/http"
	"ants/util"
	"log"
	"spiders"
)

type Crawler struct {
	SpiderMap     map[string]*util.Spider
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
	spiderMap := make(map[string]*util.Spider)
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
	if spider.Status == util.SPIDERS_STATUS_RUNNING {
		result.Success = false
		result.Message = "spider already runing"
		result.Spider = spider.Name
		return result
	}
	spider.Status = util.SPIDERS_STATUS_RUNNING
	for _, url := range spider.StartUrls {
		request, err := http.NewRequest("GET", url, nil, spider.Name, util.BASE_PARSE_NAME)
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

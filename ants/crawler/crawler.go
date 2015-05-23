package crawler

import (
	"github.com/wcong/ants-go/ants/http"
	base_spider "github.com/wcong/ants-go/ants/spiders"
	"github.com/wcong/ants-go/ants/util"
	"github.com/wcong/ants-go/spiders"
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
func NewCrawler(resultQuene *ResultQuene, settings *util.Settings) *Crawler {
	requestQuene := NewRequestQuene()
	responseQuene := NewResponseQuene()
	spiderMap := spiders.LoadAllSpiders()
	downloader := NewDownloader(requestQuene, responseQuene, settings.DownloadInterval)
	scraper := NewScraper(resultQuene, responseQuene, spiderMap)
	crawler := &Crawler{spiderMap, requestQuene, responseQuene, downloader, scraper}
	return crawler
}

func (this *Crawler) GetStartRequest(spiderName string) []*http.Request {
	spider := this.SpiderMap[spiderName]
	startRequests := spider.MakeStartRequests()
	return startRequests
}

func (this *Crawler) StartSpider(spiderName string) {
	beforeMethod := this.SpiderMap[spiderName].BeforeMethod
	if beforeMethod != nil {
		beforeMethod()
	}
}

func (this *Crawler) CloseSpider(spiderName string) {
	afterMethod := this.SpiderMap[spiderName].AfterMethod
	if afterMethod != nil {
		afterMethod()
	}
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

func (this *Crawler) Stop() {
	this.Downloader.Stop()
	this.Scraper.Stop()
}

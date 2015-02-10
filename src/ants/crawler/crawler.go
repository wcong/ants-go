package crawler

import (
	"ants/http"
	"ants/util"
	"log"
	"spiders"
)

type Crawler struct {
	SpiderMap map[string]*util.Spider
}

func (this *Crawler) LoadSpiders() {
	this.SpiderMap = make(map[string]*util.Spider)
	deadLoopTest := spiders.MakeDeadLoopSpider()
	this.SpiderMap[deadLoopTest.Name] = deadLoopTest
}
func (this *Crawler) StartSpider(spiderName string) {
	log.Println("start to crawl spider " + spiderName)
	spider := this.SpiderMap[spiderName]
	for _, url := range spider.StartUrls {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
			continue
		}
		go this.RunSpider(spider, request)
	}
}
func (this *Crawler) RunSpider(spider *util.Spider, request *http.Request) {
	response := Download(request)
	if response != nil {
		return
	}
	requestList, err := spider.ParseMap[response.ParserName](response)
	if err != nil {
		log.Fatal(err)
	}
	if requestList == nil {
		return
	}
	requestListLength := len(requestList)
	if requestListLength == 1 {
		this.RunSpider(spider, requestList[0])
	} else {
		this.RunSpider(spider, requestList[0])
		for i := 1; i < requestListLength; i++ {
			go this.RunSpider(spider, requestList[i])
		}
	}
}

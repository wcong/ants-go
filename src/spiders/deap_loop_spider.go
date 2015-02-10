package spiders

import (
	"ants/http"
	"ants/util"
	"github.com/PuerkitoBio/goquery"
)

func MakeDeadLoopSpider() *util.Spider {
	spider := util.Spider{}
	spider.Name = "deal_loop_spider"
	spider.StartUrls = []string{"http://www.baidu.com/s?wd=1"}
	spider.ParseMap = make(map[string]func(response *http.Response) ([]*http.Request, error))
	spider.ParseMap["base"] = func(response *http.Response) ([]*http.Request, error) {
		doc, err := goquery.NewDocumentFromReader(response.GoResponse.Body)
		val, _ := doc.Find("#page a[text()=\"下一页\"]").Attr("href")
		request, err := http.NewRequest("GET", val, nil)
		requestList := make([]*http.Request, 0)
		requestList = append(requestList, request)
		return requestList, err
	}
	return &spider
}

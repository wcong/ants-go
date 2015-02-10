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
	spider.Parse = func(response *http.Response) (*http.Request, error) {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		val, _ := doc.Find("#page a[text()=\"下一页\"]").Attr("href")
		request, err := http.NewRequest("GET", val, nil)
		return request, err
	}
	return &spider
}

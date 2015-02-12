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
		var nextPage string
		doc.Find("#page a").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "下一页>" {
				href, _ := s.Attr("href")
				nextPage = href
			}
		})
		nextPage = "http://www.baidu.com" + nextPage
		request, err := http.NewRequest("GET", nextPage, nil, spider.Name, util.BASE_PARSE_NAME)
		requestList := make([]*http.Request, 0)
		requestList = append(requestList, request)
		return requestList, err
	}
	return &spider
}

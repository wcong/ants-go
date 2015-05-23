package spiders

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"strings"
)

func MakeDeadLoopSpider() *spiders.Spider {
	spider := spiders.Spider{}
	spider.Name = "deal_loop_spider"
	spider.StartUrls = []string{"http://www.baidu.com/s?wd=1"}
	spider.ParseMap = make(map[string]func(response *http.Response) ([]*http.Request, error))
	spider.ParseMap[spiders.BASE_PARSE_NAME] = func(response *http.Response) ([]*http.Request, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
		if err != nil {
			return nil, err
		}
		nodes := doc.Find("#page .n").Nodes
		if len(nodes) == 0 {
			return nil, err
		}
		nextNode := nodes[len(nodes)-1]
		attrList := nextNode.Attr
		var nextPageLink string
		for _, attr := range attrList {
			if attr.Key == "href" {
				nextPageLink = attr.Val
			}
		}
		nextPage := "http://www.baidu.com" + nextPageLink
		request, err := http.NewRequest("GET", nextPage, spider.Name, spiders.BASE_PARSE_NAME, nil, 0)
		requestList := make([]*http.Request, 0)
		requestList = append(requestList, request)
		return requestList, nil
	}
	return &spider
}

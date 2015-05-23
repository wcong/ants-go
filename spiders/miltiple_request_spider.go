package spiders

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"log"
	"strings"
)

func MakeMuiltiplySpiders() *spiders.Spider {
	spider := spiders.Spider{}
	spider.Name = "muiltiply_spider"
	spider.StartUrls = []string{"http://www.baidu.com/s?wd=1"}
	spider.ParseMap = make(map[string]func(response *http.Response) ([]*http.Request, error))
	spider.ParseMap[spiders.BASE_PARSE_NAME] = func(response *http.Response) ([]*http.Request, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
		if err != nil {
			return nil, err
		}
		requestList := make([]*http.Request, 0, 10)
		doc.Find("#page a").Each(func(index int, hrefNode *goquery.Selection) {
			href, isExist := hrefNode.Attr("href")
			if !isExist {
				return
			}
			nextPage := "http://www.baidu.com" + href
			request, err := http.NewRequest("GET", nextPage, spider.Name, spiders.BASE_PARSE_NAME, nil, 0)
			if err != nil {
				log.Println(err)
			}
			requestList = append(requestList, request)
		})
		return requestList, nil
	}
	return &spider
}

package spiders

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

func Base(response *http.Response) ([]*http.Request, error) {
	requestList := make([]*http.Request, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
	if err != nil {
		return requestList, err
	}
	xsrf, exist := doc.Find(".zu-side-login-box input[name=_xsrf]").Attr("value")
	if !exist {
		return requestList, nil
	}
	userList := [][]string{{"1203316364@qq.com", ""}, {"wc19920415@163.com", ""}}
	for index, user := range userList {
		value := make(url.Values)
		value.Set("_xsrf", xsrf)
		value.Set("email", user[0])
		value.Set("password", user[1])
		value.Set("rememberme", "y")
		request, requestErr := http.NewRequest("POST", "http://www.zhihu.com/login", zhihuSpider.Name, "Index", strings.NewReader(value.Encode()), index+2)
		if requestErr != nil {
			log.Println(err)
			continue
		}
		requestList = append(requestList, request)
	}
	return requestList, nil
}

func Index(response *http.Response) ([]*http.Request, error) {
	requestList := make([]*http.Request, 1)
	request, _ := http.NewRequest("GET", "http://www.zhihu.com", zhihuSpider.Name, "GetId", nil, response.Request.CookieJar)
	requestList[0] = request
	return requestList, nil
}
func GetId(response *http.Response) ([]*http.Request, error) {
	html, _ := ioutil.ReadAll(response.GoResponse.Body)
	log.Println(string(html))
	doc, err := goquery.NewDocumentFromReader(response.GoResponse.Body)
	if err != nil {
		return nil, err
	}
	id, exist := doc.Find(".zu-top-nav-userinfo").Attr("href")
	if exist {
		log.Println(id)
	}
	return nil, nil
}

var zhihuSpider *spiders.Spider

func MakeZhiHuSpider() *spiders.Spider {
	zhihuSpider = &spiders.Spider{}
	zhihuSpider.Name = "zhihu_spider"
	zhihuSpider.StartUrls = []string{"http://www.zhihu.com/"}
	parseMap := make(map[string]func(response *http.Response) ([]*http.Request, error))
	zhihuSpider.ParseMap = parseMap
	parseMap[spiders.BASE_PARSE_NAME] = Base
	parseMap["Index"] = Index
	parseMap["GetId"] = GetId
	return zhihuSpider
}

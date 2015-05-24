package spiders

import (
	"encoding/json"
	"github.com/wcong/ants-go/ants/db"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"log"
	"math"
	"strconv"
	"strings"
)

type Stock struct {
	Code   int
	Day    string
	Count  int
	Fields []string
	Items  [][]interface{}
}

func MakeStockSpider() *spiders.Spider {
	spider := spiders.Spider{}
	spider.Name = "stock_spider"
	spider.BeforeMethod = func() {
		db.DefaultMysqlConnectMap.InitConnection(spider.Name, "root:@/wxstock?charset=utf8")
	}
	spider.AfterMethod = func() {
		db.DefaultMysqlConnectMap.CloseContection(spider.Name)
	}
	urlPrefix := "http://money.finance.sina.com.cn/d/api/openapi_proxy.php/?__s=[[%22hq%22,%22hs_a%22,%22%22,0,"
	urlSuffix := ",40]]&callback=FDC_DC.theTableData"
	pageSize := 40
	spider.StartUrls = []string{urlPrefix + "0" + urlSuffix}
	spider.ParseMap = make(map[string]func(response *http.Response) ([]*http.Request, error))
	spider.ParseMap[spiders.BASE_PARSE_NAME] = func(response *http.Response) ([]*http.Request, error) {
		body := response.Body
		startIndex := strings.Index(body, "FDC_DC.theTableData")
		jsonString := body[startIndex+len("FDC_DC.theTableData(") : len(body)-2]
		var jsonResult []Stock
		err := json.Unmarshal([]byte(jsonString), &jsonResult)
		if err != nil {
			log.Println(err)
		}
		for _, item := range jsonResult[0].Items {
			code := item[1]
			name := item[2]
			rows, err := db.DefaultMysqlConnectMap.Query(spider.Name, "select id from db_stock where code= ? ", code)
			if err != nil {
				log.Println(err)
				continue
			}
			if !rows.Next() {
				rows.Close()
				db.DefaultMysqlConnectMap.Exec(spider.Name, "insert into db_stock(gmt_create,code,name)values(now(), ?, ? )", code, name)
			} else {
				rows.Close()
			}
		}
		if strings.Contains(response.GoResponse.Request.URL.RawQuery, ",0,0,40]]") {
			totalPage := int(math.Ceil(float64(jsonResult[0].Count) / float64(pageSize)))
			requestList := make([]*http.Request, totalPage, totalPage)
			pageNo := 1
			for pageNo <= totalPage {
				request, err := http.NewRequest("GET", urlPrefix+strconv.Itoa(pageNo)+urlSuffix, spider.Name, spiders.BASE_PARSE_NAME, nil, 0)
				if err != nil {
					log.Println(err)
				}
				requestList[pageNo-1] = request
				pageNo += 1
			}
			return requestList, nil
		}
		return nil, nil
	}
	return &spider
}

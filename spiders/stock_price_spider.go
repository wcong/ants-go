package spiders

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wcong/ants-go/ants/db"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/spiders"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func MakeStockPriceSpider() *spiders.Spider {
	spider := spiders.Spider{}
	spider.Name = "stock_price_spider"
	urlPrefix := "http://money.finance.sina.com.cn/corp/go.php/vMS_MarketHistory/stockid/"
	spider.BeforeMethod = func() {
		db.DefaultMysqlConnectMap.InitConnection(spider.Name, "root:@/wxstock?charset=utf8")
		rows, err := db.DefaultMysqlConnectMap.Query(spider.Name, "select id,code from db_stock")
		if err != nil {
			log.Println(err)
		} else {
			defer rows.Close()
			spider.ExtData = make(map[string]interface{})
			codeIdMap := make(map[string]int)
			spider.ExtData["codeIdMap"] = codeIdMap
			for rows.Next() {
				var id int
				var code string
				rows.Scan(&id, &code)
				codeIdMap[code] = id
			}
		}
		priceRow, priceErr := db.DefaultMysqlConnectMap.Query(spider.Name, "select id,date from db_stock_price")
		if priceErr != nil {
			log.Println(err)
		} else {
			defer priceRow.Close()
			existMap := make(map[string]bool)
			spider.ExtData["existMap"] = existMap
			for priceRow.Next() {
				var id int
				var date string
				priceRow.Scan(&id, &date)
				existMap[strconv.Itoa(id)+":"+date] = true
			}
		}
		maxDateRow, maxErr := db.DefaultMysqlConnectMap.Query(spider.Name, "select max(date),stock_id from db_stock_price group by stock_id")
		if maxErr != nil {
			log.Println(err)
		} else {
			defer maxDateRow.Close()
			maxDataMap := make(map[int]string)
			spider.ExtData["maxDataMap"] = maxDataMap
			for maxDateRow.Next() {
				var maxDate string
				var stockId int
				maxDateRow.Scan(&maxDate, &stockId)
				maxDataMap[stockId] = maxDate
			}
		}
	}
	spider.AfterMethod = func() {
		db.DefaultMysqlConnectMap.CloseContection(spider.Name)
	}
	spider.StartUrls = []string{"http://money.finance.sina.com.cn/corp/go.php/vMS_MarketHistory/stockid/600570.phtml?year=2015&jidu=2"}
	spider.ParseMap = make(map[string]func(response *http.Response) ([]*http.Request, error))
	spider.ParseMap[spiders.BASE_PARSE_NAME] = func(response *http.Response) ([]*http.Request, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		code := response.Request.GoRequest.URL.Path[len("/corp/go.php/vMS_MarketHistory/stockid/"):len("/corp/go.php/vMS_MarketHistory/stockid/600570")]
		doc.Find("#FundHoldSharesTable tr").Each(func(index int, selection *goquery.Selection) {
			if index == 0 {
				return
			}
			date := strings.TrimSpace(selection.First().Find("div a").Text())
			if len(date) == 0 {
				return
			}
			priceString := strings.TrimSpace(selection.Find("td:nth-child(3) div").Text())
			priceFloat, _ := strconv.ParseFloat(priceString, 32)
			price := int(priceFloat * 1000)
			id := (spider.ExtData["codeIdMap"]).(map[string]int)[code]
			_, ok := (spider.ExtData["existMap"]).(map[string]bool)[strconv.Itoa(id)+":"+date]
			if !ok {
				_, err := db.DefaultMysqlConnectMap.Exec(spider.Name, "insert into db_stock_price(gmt_create,creator,stock_id,date,price)values(now(),'go', ? , ? , ? )", id, date, price)
				if err != nil {
					log.Println(err)
				}
			}
		})
		return nil, nil
	}
	spider.InitStartUrls = func() {
		spider.StartUrls = make([]string, 0)
		today := time.Now()
		year := today.Year()
		month := today.Month()
		jidu := int(math.Ceil(float64(month) / float64(3)))
		codeIdMap := (spider.ExtData["codeIdMap"]).(map[string]int)
		maxDataMap := (spider.ExtData["maxDataMap"]).(map[int]string)
		for code, id := range codeIdMap {
			stringData, ok := maxDataMap[id]
			var maxYear int
			var maxJidu int
			if ok {
				maxDate, _ := time.Parse("2006-01-02", stringData)
				maxYear = maxDate.Year()
				maxMonth := maxDate.Month()
				maxJidu = int(math.Ceil(float64(maxMonth) / float64(3)))
			} else {
				maxYear = year
				maxJidu = jidu
			}
			var i int = 1
			if maxYear >= year {
				i = maxJidu
				for ; i <= jidu; i++ {
					spider.StartUrls = append(spider.StartUrls, urlPrefix+code+".phtml?year="+strconv.Itoa(year)+"&jidu="+strconv.Itoa(i))
				}
			} else {
				for ; i <= jidu; i++ {
					spider.StartUrls = append(spider.StartUrls, urlPrefix+code+".phtml?year="+strconv.Itoa(year)+"&jidu="+strconv.Itoa(i))
				}
				if i < maxJidu {
					i = maxJidu
				}
				for ; i < 5; i++ {
					spider.StartUrls = append(spider.StartUrls, urlPrefix+code+".phtml?year="+strconv.Itoa(year-1)+"&jidu="+strconv.Itoa(i))
				}
			}
		}
	}
	return &spider
}

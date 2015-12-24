package spiders

// where i defined spider
import (
	"github.com/wcong/ants-go/ants/http"
	"log"
)

const (
	BASE_PARSE_NAME     = "base"
	SPIDERS_STATUS_INIT = iota
	SPIDERS_STATUS_RUNNING
	SPIDERS_STATUS_STOP
	SPIDERS_BASIC_COOKIE
)

/*
what a spider do
*	make start request
* 	define basic parse func
*/
type Spider struct {
	BeforeMethod  func()
	InitStartUrls func()
	AfterMethod   func()
	Name          string
	StartUrls     []string
	ExtData       map[string]interface{}
	ParseMap      map[string]func(response *http.Response) ([]*http.Request, error) // defined you own parse function to scrapy response,make sure it return []response
}

func (this *Spider) MakeStartRequests() []*http.Request {
	if this.InitStartUrls != nil {
		this.InitStartUrls()
	}
	startRequestSlice := make([]*http.Request, len(this.StartUrls))
	for index, url := range this.StartUrls {
		request, err := http.NewRequest("GET", url, this.Name, BASE_PARSE_NAME, nil, 0)
		if err != nil {
			log.Println(err)
			continue
		}
		startRequestSlice[index] = request
	}
	return startRequestSlice
}
### how to write a spider
#### basic concept
*	request
	*	CookieJar int: if the web site store message in cookie which show diffrent result,you shoud use it
	*	ParserName string: crawler will send the response to named parseFunction
*	ParseMap, defined you own parse function to scrapy response,make sure it return []response

#### request struct
```go
type Spider struct {
	Name      string // unique name of spider,start a spider by name
	StartUrls []string // start urls
	ParseMap  map[string]func(response *http.Response) ([]*http.Request, error) // parse map ,we explain it above
}
```
#### parse function
**use goquery parse html**
**for example we parse a url from baidu**

*	accept response as param
*	return request slice and error if exist
*	parse html elements by go query,get next request url

```go
	parseMap["base"] = func(response *http.Response) ([]*http.Request, error) {
		if response.Request.Depth > 10 {
			return nil, nil
		}
		doc, err := goquery.NewDocumentFromReader(response.GoResponse.Body)
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
				break
			}
		}
		nextPage := "http://www.baidu.com" + nextPageLink
		request, err := http.NewRequest("GET", nextPage, "spider_name", "parse_name", nil, 0)
		requestList := make([]*http.Request, 0)
		requestList = append(requestList, request)
		return requestList, nil
	}
```

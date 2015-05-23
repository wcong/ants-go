package http

import (
	"io/ioutil"
	"log"
	Http "net/http"
)

const (
	CONTENT_TYPE      = "Content-Type"
	JSON_CONTENT_TYPE = "application/json"
)

type Response struct {
	GoResponse *Http.Response
	Request    *Request
	SpiderName string
	ParserName string
	NodeName   string
	Body       string
}

func NewResponse(response *Http.Response, request *Request, spiderName, parserName, nodeName string) *Response {
	var body []byte
	var err error
	if response != nil {
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
		}
		response.Body.Close()
	} else {
		body = make([]byte, 0, 0)
	}
	return &Response{response, request, spiderName, parserName, nodeName, string(body)}
}

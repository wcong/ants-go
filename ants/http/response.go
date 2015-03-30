package http

import (
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
}

func NewResponse(response *Http.Response, request *Request, spiderName, parserName, nodeName string) *Response {
	return &Response{response, request, spiderName, parserName, nodeName}
}

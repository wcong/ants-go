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
	SpiderName string
	ParserName string
	NodeName   string
}

func NewResponse(response *Http.Response, spiderName, parserName, nodeName string) *Response {
	return &Response{response, spiderName, parserName, nodeName}
}

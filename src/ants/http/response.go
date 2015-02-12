package http

import (
	Http "net/http"
)

type Response struct {
	GoResponse *Http.Response
	SpiderName string
	ParserName string
}

func NewResponse(response *Http.Response, spiderName string, parserName string) *Response {
	return &Response{response, spiderName, parserName}
}

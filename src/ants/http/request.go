package http

import (
	"fmt"
	"io"
	Http "net/http"
)

type Request struct {
	GoRequest  *Http.Request
	SpiderName string
	ParserName string
}

func NewRequest(method, url string, body io.Reader, spiderName string, parserName string) (*Request, error) {
	httpRequest, err := Http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request := Request{httpRequest, spiderName, parserName}
	return &request, err
}

func PrintPackage() {
	fmt.Println("this is request package")
}

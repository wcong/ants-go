package http

import (
	"fmt"
	"io"
	Http "net/http"
)

type Request struct {
	Http.Request
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	httpRequest, err := Http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request := Request{*httpRequest}
	return &request, err
}

func PrintPackage() {
	fmt.Println("this is request package")
}

package http

import (
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	http.Request
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	httpRequest, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return httpRequest, err
}

func PrintPackage() {
	fmt.Println("this is request package")
}

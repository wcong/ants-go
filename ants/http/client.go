package http

import (
	Http "net/http"
	"net/url"
	"time"
)

var proxyMap map[string]*Http.Transport = make(map[string]*Http.Transport)

type Client struct {
	GoClient *Http.Client
}

func NewClient() *Client {
	goClient := &Http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{goClient}
}
func (this *Client) SetProxy(urlString string) {
	transport, ok := proxyMap[urlString]
	if !ok {
		proxyUrl, _ := url.Parse(urlString)
		transport = &Http.Transport{Proxy: Http.ProxyURL(proxyUrl)}
		proxyMap[urlString] = transport
	}
	this.GoClient.Transport = transport
}
func (this *Client) ClearProxy() {
	this.GoClient.Transport = nil
}

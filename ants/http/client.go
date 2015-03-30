package http

import (
	Http "net/http"
)

type Client struct {
	GoClient *Http.Client
}

func NewClient() *Client {
	goClient := &Http.Client{}
	return &Client{goClient}
}

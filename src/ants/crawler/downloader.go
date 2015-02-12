package crawler

import (
	"ants/http"
	"log"
	"time"
)

const (
	DOWNLOADER_STATUS_STOP = iota
	DOWNLOADER_STATUS_RUNING
)

type Downloader struct {
	Status        int
	RequestQuene  *RequestQuene
	ResponseQuene *ResponseQuene
	ClientList    []*http.Client
}

func NewDownloader(resuqstQuene *RequestQuene, responseQuene *ResponseQuene) *Downloader {
	clientList := make([]*http.Client, 0)
	client := http.NewClient()
	clientList = append(clientList, client)
	return &Downloader{DOWNLOADER_STATUS_STOP, resuqstQuene, responseQuene, clientList}
}

func (this *Downloader) Start() {
	if this.Status == DOWNLOADER_STATUS_RUNING {
		return
	}
	log.Println("start downloader")
	this.Status = DOWNLOADER_STATUS_RUNING
	this.Download()
}
func (this *Downloader) Stop() {
	this.Status = DOWNLOADER_STATUS_STOP
}
func (this *Downloader) Download() {
	for {
		if this.Status != DOWNLOADER_STATUS_RUNING {
			break
		}
		request := this.RequestQuene.Pop()
		if request == nil {
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("download url:" + request.GoRequest.URL.String())
		client := this.ClientList[0]
		response, err := client.GoClient.Do(request.GoRequest)
		if err != nil {
			log.Fatal(err)
		}
		Response := http.NewResponse(response, request.SpiderName, request.ParserName)
		this.ResponseQuene.Push(Response)
	}
}

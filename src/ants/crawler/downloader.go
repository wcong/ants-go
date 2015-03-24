package crawler

import (
	"ants/http"
	"log"
	"time"
)

const (
	DOWNLOADER_STATUS_STOP = iota
	DOWNLOADER_STATUS_RUNING
	DOWNLOADER_STATUS_PAUSE
	DOWNLOADER_STATUS_STOPED
)

// downloader tools
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
	return &Downloader{DOWNLOADER_STATUS_STOPED, resuqstQuene, responseQuene, clientList}
}

// DOWNLOADER_STATUS_STOPED means the dead loop is actually dead
func (this *Downloader) Start() {
	if this.Status == DOWNLOADER_STATUS_RUNING {
		return
	}
	for {
		if this.IsStop() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("start downloader")
	this.Status = DOWNLOADER_STATUS_RUNING
	this.Download()
}
func (this *Downloader) Stop() {
	this.Status = DOWNLOADER_STATUS_STOP
}
func (this *Downloader) Pause() {
	if this.Status == DOWNLOADER_STATUS_RUNING {
		this.Status = DOWNLOADER_STATUS_PAUSE
	}
}
func (this *Downloader) UnPause() {
	if this.Status == DOWNLOADER_STATUS_PAUSE {
		this.Status = DOWNLOADER_STATUS_RUNING
	}
}
func (this *Downloader) IsStop() bool {
	return this.Status == DOWNLOADER_STATUS_STOPED
}

// dead loop for download
// pop a request
// download it
// push to response quene
func (this *Downloader) Download() {
	for {
		if this.Status == DOWNLOADER_STATUS_PAUSE {
			time.Sleep(1 * time.Second)
			continue
		}
		if this.Status == DOWNLOADER_STATUS_STOP {
			this.Status = DOWNLOADER_STATUS_STOPED
			break
		}
		request := this.RequestQuene.Pop()
		if request == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		go this.downloadAndPush(request)
	}
}

// download it and push in goroutine
func (this *Downloader) downloadAndPush(request *http.Request) {
	log.Println(request.SpiderName, "depth:", request.Depth, "download url:", request.GoRequest.URL.String())
	client := this.ClientList[0]
	response, err := client.GoClient.Do(request.GoRequest)
	if err != nil {
		log.Println(err)
	}
	Response := http.NewResponse(response, request, request.SpiderName, request.ParserName, request.NodeName)
	this.ResponseQuene.Push(Response)
}

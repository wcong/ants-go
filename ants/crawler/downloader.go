package crawler

import (
	"github.com/wcong/ants-go/ants/http"
	"log"
	"net/http/cookiejar"
	"sync"
	"time"
)

var cookieMutex sync.Mutex

const (
	DOWNLOADER_STATUS_STOP = iota
	DOWNLOADER_STATUS_RUNING
	DOWNLOADER_STATUS_PAUSE
	DOWNLOADER_STATUS_STOPED
)

// downloader tools
type Downloader struct {
	Status           int
	RequestQuene     *RequestQuene
	ResponseQuene    *ResponseQuene
	ClientList       []*http.Client
	DownloadInterval int
}

func NewDownloader(resuqstQuene *RequestQuene, responseQuene *ResponseQuene, downloadInterval int) *Downloader {
	clientList := make([]*http.Client, 0)
	client := http.NewClient()
	clientList = append(clientList, client)
	return &Downloader{DOWNLOADER_STATUS_STOPED, resuqstQuene, responseQuene, clientList, downloadInterval}
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
	clientList := make([]*http.Client, 0)
	client := http.NewClient()
	clientList = append(clientList, client)
	this.ClientList = clientList
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
		if this.DownloadInterval > 0 {
			time.Sleep(time.Duration(this.DownloadInterval) * time.Second)
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
	client := this.getClient(request)
	if request.Proxy != "" {
		client.SetProxy(request.Proxy)
	}
	response, err := client.GoClient.Do(request.GoRequest)
	if err != nil {
		log.Println(err)
		if request.Retry < 3 {
			log.Println(request.SpiderName, "error time", request.Retry, "return to quene:", request.GoRequest.URL.String())
			request.Retry += 1
			this.RequestQuene.Push(request)
			return
		} else {
			log.Println(request.SpiderName, "error time", request.Retry, "drop:", request.GoRequest.URL.String())
		}
	}
	client.ClearProxy()
	Response := http.NewResponse(response, request, request.SpiderName, request.ParserName, request.NodeName)
	this.ResponseQuene.Push(Response)
}

// get client by cookie jar index
func (this *Downloader) getClient(request *http.Request) *http.Client {
	this.makeClientIfNotExist(request.CookieJar)
	return this.ClientList[request.CookieJar]
}

// make cookie jar if not exist
func (this *Downloader) makeClientIfNotExist(cookieJar int) {
	clientLength := len(this.ClientList)
	if (clientLength - 1) >= cookieJar {
		return
	}
	cookieMutex.Lock()
	fixEnd := cookieJar + 1
	log.Println("none cookie jar", cookieJar, "so make up to it")
	for i := clientLength; i < fixEnd; i++ {
		jar, _ := cookiejar.New(nil)
		client := http.NewClient()
		client.GoClient.Jar = jar
		this.ClientList = append(this.ClientList, client)
	}
	cookieMutex.Unlock()
}

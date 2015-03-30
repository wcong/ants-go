package crawler

import (
	"container/list"
	"github.com/wcong/ants-go/ants/http"
	"sync"
)

var mutex sync.Mutex

type RequestQuene struct {
	RequestList *list.List
}

func NewRequestQuene() *RequestQuene {
	requestList := RequestQuene{}
	requestList.RequestList = list.New()
	return &requestList
}

func (this *RequestQuene) Push(request *http.Request) {
	this.RequestList.PushBack(request)
}
func (this *RequestQuene) Pop() *http.Request {
	element := this.RequestList.Front()
	if element == nil {
		return nil
	}
	this.RequestList.Remove(element)
	value := element.Value
	return value.(*http.Request)
}

func (this *RequestQuene) IsEmpty() bool {
	return this.RequestList.Len() == 0
}

type ResponseQuene struct {
	ResponseList *list.List
}

func NewResponseQuene() *ResponseQuene {
	requestList := ResponseQuene{}
	requestList.ResponseList = list.New()
	return &requestList
}

func (this *ResponseQuene) Push(response *http.Response) {
	mutex.Lock()
	this.ResponseList.PushBack(response)
	mutex.Unlock()
}

// for now only one routine pop the request,so do not add lock
func (this *ResponseQuene) Pop() *http.Response {
	element := this.ResponseList.Front()
	if element == nil {
		return nil
	}
	this.ResponseList.Remove(element)
	value := element.Value
	return value.(*http.Response)
}

type ResultQuene struct {
	ResultList *list.List
}

func NewResultQuene() *ResultQuene {
	resultList := list.New()
	resultQuene := &ResultQuene{resultList}
	return resultQuene
}
func (this *ResultQuene) Push(scrapeResult *ScrapeResult) {
	mutex.Lock()
	this.ResultList.PushBack(scrapeResult)
	mutex.Unlock()
}
func (this *ResultQuene) Pop() *ScrapeResult {
	element := this.ResultList.Front()
	if element == nil {
		return nil
	}
	this.ResultList.Remove(element)
	value := element.Value
	return value.(*ScrapeResult)
}

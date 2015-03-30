package http

import (
	"github.com/wcong/ants-go/ants/util"
	"log"
	Http "net/http"
	"strconv"
	"sync"
)

type HttpServer struct {
	Http.Server
}

func NewHttpServer(setting *util.Settings, handler Http.Handler) *HttpServer {
	port := strconv.Itoa(setting.HttpPort)
	httpServer := &HttpServer{
		Http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
	return httpServer
}

func (this *HttpServer) Start(wg sync.WaitGroup) {
	go this.server(wg)
}

func (this *HttpServer) server(wg sync.WaitGroup) {
	log.Println("start to server http" + this.Addr)
	err := this.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("http server down")
	wg.Done()
}

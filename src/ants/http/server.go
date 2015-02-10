package http

import (
	"ants/conf"
	"log"
	Http "net/http"
	"strconv"
	"sync"
)

type HttpServer struct {
	Http.Server
}

func NewHttpServer(setting *conf.Settings, handler Http.Handler) *HttpServer {
	port := strconv.Itoa(setting.HttpPort)
	httpServer := &HttpServer{
		Http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
	return httpServer
}

func (this *HttpServer) Start(wg *sync.WaitGroup) {
	log.Println("start to server http" + this.Addr)
	err := this.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
	wg.Done()
}

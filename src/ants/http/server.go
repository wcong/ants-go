package http

import (
	"ants/http/server"
	"log"
	Http "net/http"
)

type HttpServer struct {
	Http.Server
}

func InitServer(ch chan int) *HttpServer {
	log.Println("going to start  server http")
	serveMux := Http.NewServeMux()
	serveMux.HandleFunc("/", server.Welcome)
	serveMux.HandleFunc("/spiders", server.Spiders)
	httpServer := &HttpServer{
		Http.Server{
			Addr:    ":8200",
			Handler: serveMux,
		},
	}

	log.Println("going to start  server http")
	go func() {
		err := httpServer.ListenAndServe()
		log.Println("start to server http")
		ch <- 1
		if err != nil {
			log.Panicln(err)
		}
	}()
	return httpServer
}

package crawler

import (
	"ants/http"
	"log"
	Http "net/http"
)

func Download(request *http.Request) *http.Response {
	log.Println("download url:" + request.GoRequest.URL.String())
	client := Http.Client{}
	response, err := client.Do(request.GoRequest)
	if err != nil {
		log.Fatal(err)
	}
	Response := http.Response{}
	Response.GoResponse = response
	return &Response
}

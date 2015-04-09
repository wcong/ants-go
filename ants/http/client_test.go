package http

import (
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient()
	request, _ := NewRequest("GET", "http://www.baidu.com", "", "", nil, 0)
	client.SetProxy("http://217.24.252.250:3128")
	response, err := client.GoClient.Do(request.GoRequest)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.Status)
	}
}

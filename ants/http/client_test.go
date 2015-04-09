package http

import (
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient()
	client.SetProxy("http://217.24.252.250:3128")
	testGetResponse(client)
}

func TestClearProxy(t *testing.T) {
	client := NewClient()
	client.SetProxy("http://217.24.252.250:3128")
	client.ClearProxy()
	testGetResponse(client)
}
func testGetResponse(client *Client) {
	request, _ := NewRequest("GET", "http://www.baidu.com", "", "", nil, 0)
	response, err := client.GoClient.Do(request.GoRequest)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.Status)
	}
}

package http

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	request, _ := NewRequest("GET", "http://www.baidu.com", "a", "a", nil, 1)
	message, _ := json.Marshal(request)
	msg := string(message)
	//fmt.Println(msg)
	parsedRequest := &Request{}
	json.Unmarshal([]byte(msg), parsedRequest)
	fmt.Println(parsedRequest.GoRequest.Host)
}

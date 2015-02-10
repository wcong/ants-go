package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	ch := make(chan int)
	go InitServer(ch)
	for {
		time.Sleep(1)
		fmt.Println("loop")
		a := <-ch
		if a == 1 {
			break
		}
	}
	response, err := http.Get("http://localhost:8200/")
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()
	result, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(result))
}

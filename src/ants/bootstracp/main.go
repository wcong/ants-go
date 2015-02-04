package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func parse(url string) (request string) {

}
func main() {
	ch := make(chan int, 1)
	fmt.Println("you suck")
	response, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
		return
	}
	a := <-ch
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

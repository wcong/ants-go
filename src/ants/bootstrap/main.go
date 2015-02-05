package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func main() {
	fmt.Println("you suck")
	response, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(response)
	html, err := doc.Find("#s_tab").Html()
	fmt.Println(html)
}

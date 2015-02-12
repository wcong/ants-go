package crawler

import (
	"ants/http"
	"log"
	"testing"
)

func TestDownload(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	log.Println(request)
}

package crawler

import (
	"ants/http"
	"testing"
)

func TestDownload(t *testing.T) {
	request := http.NewRequest("GET", "http://www.baidu.com", nil)

}

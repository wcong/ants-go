package util

// where i defined spider
import (
	"ants/http"
)

const (
	BASE_PARSE_NAME = "base"
)

type Spider struct {
	Name      string
	StartUrls []string
	ParseMap  map[string]func(response *http.Response) ([]*http.Request, error)
}

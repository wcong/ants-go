package util

// where i defined spider
import (
	"ants/http"
)

type Spider struct {
	Name      string
	StartUrls []string
	ParseMap  map[string]func(response *http.Response) ([]*http.Request, error)
}

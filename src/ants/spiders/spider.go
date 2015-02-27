package spiders

// where i defined spider
import (
	"ants/http"
)

const (
	BASE_PARSE_NAME     = "base"
	SPIDERS_STATUS_INIT = iota
	SPIDERS_STATUS_RUNNING
	SPIDERS_STATUS_STOP
)

type Spider struct {
	Status    int
	Name      string
	StartUrls []string
	ParseMap  map[string]func(response *http.Response) ([]*http.Request, error)
}

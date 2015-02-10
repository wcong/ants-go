package http

import (
	Http "net/http"
)

type Response struct {
	GoResponse *Http.Response
	ParserName string "base"
}

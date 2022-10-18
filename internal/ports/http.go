package ports

import (
	"net/http"
)

type HttpServer interface {
	Serve() error
}

type Route struct {
	Method  string
	Route   string
	Handler http.Handler
}

package ports

import "net/http"

type HandlerProvider interface {
	GetHandler() http.Handler
}

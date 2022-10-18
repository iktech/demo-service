package http

import (
	"github.com/iktech/demo-service/internal/ports"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ServerBuilder struct {
	server *Server
}

func NewServerBuilder(address *net.IPAddr, port int) *ServerBuilder {
	return &ServerBuilder{
		server: &Server{
			ipAddress:    address.String(),
			port:         strconv.Itoa(port),
			readTimeout:  5,
			writeTimeout: 10,
			idleTimeout:  15,
			routes:       make([]ports.Route, 0),
		},
	}
}

func (b *ServerBuilder) WithReadTimeout(readTimeout int) *ServerBuilder {
	b.server.readTimeout = time.Duration(readTimeout)

	return b
}

func (b *ServerBuilder) WithWriteTimeout(readTimeout int) *ServerBuilder {
	b.server.readTimeout = time.Duration(readTimeout)

	return b
}

func (b *ServerBuilder) IdleWriteTimeout(readTimeout int) *ServerBuilder {
	b.server.readTimeout = time.Duration(readTimeout)

	return b
}

func (b *ServerBuilder) AddRoute(method, path string, handler http.Handler) *ServerBuilder {
	b.server.routes = append(b.server.routes, ports.Route{
		Method:  method,
		Route:   path,
		Handler: handler,
	})

	return b
}

func (b *ServerBuilder) Build() *Server {
	return b.server
}

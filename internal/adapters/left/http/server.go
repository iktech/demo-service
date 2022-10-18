package http

import (
	"context"
	"github.com/iktech/demo-service/internal/ports"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	ipAddress    string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
	routes       []ports.Route
}

var quit chan os.Signal

func (s Server) Serve() error {
	httpServer := &http.Server{
		Addr:         s.ipAddress + ":" + s.port,
		ReadTimeout:  s.readTimeout * time.Second,
		WriteTimeout: s.writeTimeout * time.Second,
		IdleTimeout:  s.idleTimeout * time.Second,
	}

	done := make(chan bool)
	quit = make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)

	h := httprouter.New()
	h.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
			header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Transaction-Id, Authorization")
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	httpServer.Handler = h

	if s.routes != nil {
		for _, r := range s.routes {
			h.Handler(r.Method, r.Route, r.Handler)
		}
	}

	go func() {
		<-quit
		log.Println("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpServer.SetKeepAlivesEnabled(false)
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the HTTP server: %v\n", err)
		}

		close(done)
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not start HTTP server on %s: %v\n", s.port, err)
	}
	log.Println("the HTTP server is ready to handle requests at: ", s.port)

	<-done
	log.Println("server stopped")
	return nil
}

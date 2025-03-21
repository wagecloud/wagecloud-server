package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wagecloud/wagecloud-server/internal/http/handler"
	"github.com/wagecloud/wagecloud-server/internal/service"
)

// Server represents the HTTP server
type Server struct {
	server  *http.Server
	handler *handler.Handler
}

// NewServer creates a new HTTP server
func NewServer(service *service.Service, addr string) *Server {
	handler := handler.NewHandler(service)

	return &Server{
		server: &http.Server{
			Addr:    "localhost:8080",
			Handler: handler.SetupRoutes(),
		},
		handler: handler,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server listening on %s", s.server.Addr)
		serverErrors <- s.server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select
	select {
	case err := <-serverErrors:
		return err

	case <-shutdown:
		log.Println("Server is shutting down...")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("Could not gracefully shutdown the server: %v\n", err)
			return err
		}

		log.Println("Server gracefully stopped")
	}

	return nil
}

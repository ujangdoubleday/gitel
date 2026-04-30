package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// server wraps an http.Server with graceful shutdown capabilities.
type Server struct {
	httpServer *http.Server
}

// newServer creates a new HTTP server with the given port and handler.
func NewServer(port string, webhookHandler *WebhookHandler) *Server {
	mux := http.NewServeMux()

	// go 1.22+ enhanced routing: method + path pattern.
	mux.HandleFunc("POST /github-webhook", webhookHandler.HandleGitHubWebhook)

	// health check endpoint.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// start begins listening for incoming HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// shutdown gracefully stops the server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/kube"
)

// ServerConfig holds configuration for the web server
type ServerConfig struct {
	Port          int
	CacheInterval time.Duration
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:          3000,
		CacheInterval: 20 * time.Second,
	}
}

// RunServer starts the HTTP server with all handlers and middleware
func RunServer(ctx context.Context, clients *kube.Clients, cfg ServerConfig) error {
	// Create handler with background caching
	handler := NewHandler(clients, config.GetConfig, cfg.CacheInterval)
	handler.StartBackgroundCache(ctx)

	// Create router
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/apps", handler.AppsHandler)
	mux.HandleFunc("GET /api/config", handler.ConfigHandler)

	// Health endpoints
	mux.HandleFunc("GET /healthz", handler.HealthzHandler)
	mux.HandleFunc("GET /readyz", handler.ReadyzHandler)

	// Serve frontend static files
	frontendFS := GetFrontendFS()
	fileServer := http.FileServer(frontendFS)

	// SPA fallback: serve index.html for routes that don't match static files
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly first
		path := r.URL.Path

		// For API routes, let them 404 normally
		if len(path) >= 4 && path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}

		// Try to open the file
		f, err := frontendFS.Open(path[1:]) // Remove leading slash
		if err != nil {
			// File doesn't exist, serve index.html for SPA routing
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		// File exists, serve it
		fileServer.ServeHTTP(w, r)
	})

	// Apply middleware stack
	wrapped := ChainMiddleware(mux,
		LoggingMiddleware,
		SecurityHeadersMiddleware,
		CORSMiddleware,
		CacheControlMiddleware,
		GzipMiddleware,
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      wrapped,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	logger.Info("Starting server on port ", cfg.Port)
	return server.ListenAndServe()
}

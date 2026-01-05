package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/kube"
)

// ServerConfig holds configuration for the web server
type ServerConfig struct {
	Port          int
	CacheInterval time.Duration
	BasePath      string
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:          3000,
		CacheInterval: 20 * time.Second,
		BasePath:      "",
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

	indexHTML, err := readIndexHTML(frontendFS)
	if err != nil {
		return fmt.Errorf("failed to read index.html: %w", err)
	}

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
			serveIndexWithBasePath(w, r, indexHTML)
			return
		}
		_ = f.Close()

		if path == "/" || path == "/index.html" {
			serveIndexWithBasePath(w, r, indexHTML)
			return
		}

		// File exists, serve it
		fileServer.ServeHTTP(w, r)
	})

	// Apply middleware stack
	wrapped := ChainMiddleware(mux,
		BasePathMiddleware(cfg.BasePath),
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
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Error during server shutdown: ", err)
		}
	}()

	logger.Info("Starting server on port ", cfg.Port)
	return server.ListenAndServe()
}

// readIndexHTML reads the index.html template from the frontend filesystem
func readIndexHTML(fs http.FileSystem) ([]byte, error) {
	f, err := fs.Open("index.html")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// serveIndexWithBasePath serves index.html with the basePath injected as a JavaScript variable.
// This allows the frontend to know what base path it's being served from for API calls.
func serveIndexWithBasePath(w http.ResponseWriter, r *http.Request, indexHTML []byte) {
	basePath := GetBasePath(r)

	basePathScript := fmt.Sprintf(`<script>window.__BASE_PATH__ = %q;</script></head>`, basePath)

	modified := bytes.Replace(indexHTML, []byte("</head>"), []byte(basePathScript), 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(modified)
}

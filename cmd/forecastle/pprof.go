package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* handlers on http.DefaultServeMux
	"time"
)

// startPprofServer starts a dedicated HTTP server exposing the net/http/pprof
// endpoints on the given port. It runs on its own port (never the main
// application server) so the profiling endpoints are not exposed publicly and
// can be reached via `kubectl port-forward`. The server shuts down when ctx is
// cancelled. A nil Handler uses http.DefaultServeMux, where the blank
// net/http/pprof import has registered the /debug/pprof/ routes.
func startPprofServer(ctx context.Context, port int) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Error during pprof server shutdown: ", err)
		}
	}()

	go func() {
		logger.Infof("Starting pprof server on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("pprof server error: %v", err)
		}
	}()
}

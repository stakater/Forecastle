package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/stakater/Forecastle/v1/internal/web"
	"github.com/stakater/Forecastle/v1/pkg/kube"
	"github.com/stakater/Forecastle/v1/pkg/log"
)

var logger = log.New()

func init() {
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.AddConfigPath("/etc/forecastle/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.forecastle") // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(errors.New("Fatal error config file: " + err.Error()))
	}
}

func main() {
	// Parse command line flags
	port := flag.Int("port", 3000, "Server port")
	cacheInterval := flag.Duration("cache-interval", 20*time.Second, "Background cache refresh interval")
	flag.Parse()

	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Infof("Received signal %v, initiating shutdown...", sig)
		cancel()
	}()

	// Initialize Kubernetes clients
	clients := kube.GetClients()

	// Configure server
	cfg := web.ServerConfig{
		Port:          *port,
		CacheInterval: *cacheInterval,
	}

	// Start server
	logger.Info("Forecastle starting...")
	if err := web.RunServer(ctx, &clients, cfg); err != nil {
		if err.Error() != "http: Server closed" {
			logger.Error("Server error: ", err)
			os.Exit(1)
		}
	}

	logger.Info("Forecastle stopped")
}

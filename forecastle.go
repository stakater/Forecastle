package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stakater/Forecastle/v1/pkg/handlers"
	"github.com/stakater/Forecastle/v1/pkg/log"
)

var (
	logger = log.New()
)

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
	basePath := ""
	router := mux.NewRouter()

	// API routes
	sub := router.PathPrefix(basePath).Subrouter()
	sub.HandleFunc("/api/apps", handlers.AppsHandler).Methods("GET")
	sub.HandleFunc("/api/config", handlers.ConfigHandler).Methods("GET")

	// Packr boxes for frontend
	staticBox := packr.New("static", "./frontend/build/static")
	buildBox := packr.New("build", "./frontend/build")

	// Serve static assets
	sub.PathPrefix("/static/").Handler(
		http.StripPrefix(basePath+"/static/", http.FileServer(staticBox)),
	)

	// SPA fallback handler
	sub.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip basePath from URL
		filePath := strings.TrimPrefix(r.URL.Path, basePath+"/")

		// If the path is empty (user requested /forecastle), serve index.html
		if filePath == "" {
			indexData, err := buildBox.Find("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			w.Write(indexData)
			return
		}

		// Try to find the file in buildBox
		if fileData, err := buildBox.Find(filePath); err == nil {
			w.Write(fileData)
			return
		}

		// Fallback to index.html for React routing
		indexData, err := buildBox.Find("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		w.Write(indexData)
	})

	logger.Infof("Listening at port 3000 with base path: %s", basePath)
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Error(err)
	}
}

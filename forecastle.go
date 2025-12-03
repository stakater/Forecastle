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

func writeOr500(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		logger.Error(err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/forecastle/")
	viper.AddConfigPath("$HOME/.forecastle")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
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
		filePath := strings.TrimPrefix(r.URL.Path, basePath+"/")

		// Default to index.html
		if filePath == "" {
			indexData, err := buildBox.Find("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			writeOr500(w, indexData)
			return
		}

		// Serve static file if exists
		if fileData, err := buildBox.Find(filePath); err == nil {
			writeOr500(w, fileData)
			return
		}

		// Fallback to SPA index.html
		indexData, err := buildBox.Find("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		writeOr500(w, indexData)
	})

	logger.Infof("Listening at port 3000 with base path: %s", basePath)
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Error(err)
	}
}

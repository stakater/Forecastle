package main

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/stakater/Forecastle/pkg/handlers"
	"github.com/stakater/Forecastle/pkg/log"
)

var (
	logger = log.New()
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/apps", handlers.AppsHandler)
	router.HandleFunc("/apps/", handlers.AppsHandler)

	router.Path("/apps").Queries("namespaces", "{namespaces}").HandlerFunc(handlers.AppsHandler)
	router.Path("/apps/").Queries("namespaces", "{namespaces}").HandlerFunc(handlers.AppsHandler)

	router.Path("/file").Queries("path", "{path}").HandlerFunc(handlers.FileHandler)
	router.Path("/file/").Queries("path", "{path}").HandlerFunc(handlers.FileHandler)

	// Serve static files using packr
	box := packr.NewBox("./static")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(box)))

	logger.Info("Listening at port 3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Error(err)
	}
}

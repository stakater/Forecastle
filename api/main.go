package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stakater/Forecastle/api/pkg/handlers"
	"github.com/stakater/Forecastle/api/pkg/log"
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

	logger.Info("Listening at port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		logger.Error(err)
	}
}

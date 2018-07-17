package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stakater/Forecastle/api/pkg/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/apps", handlers.AppsHandler)
	router.HandleFunc("/apps/", handlers.AppsHandler)

	router.Path("/apps").Queries("namespaces", "{namespaces}").HandlerFunc(handlers.AppsHandler)
	router.Path("/apps/").Queries("namespaces", "{namespaces}").HandlerFunc(handlers.AppsHandler)

	log.Info("Listening at port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Error(err)
	}
}

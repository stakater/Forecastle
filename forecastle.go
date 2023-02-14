package main

import (
	"errors"
	"net/http"

	packr "github.com/gobuffalo/packr/v2"
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
	router := mux.NewRouter()

	router.HandleFunc("/api/apps", handlers.AppsHandler).Methods("GET")
	router.HandleFunc("/api/config", handlers.ConfigHandler).Methods("GET")

	// Serve react frontend using packr
	staticBox := packr.New("static", "./frontend/build/static")
	buildBox := packr.New("build", "./frontend/build")
	router.PathPrefix("/").Handler(http.FileServer(buildBox))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox)))

	logger.Info("Listening at port 3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Error(err)
	}
}

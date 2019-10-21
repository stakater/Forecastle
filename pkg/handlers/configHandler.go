package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stakater/Forecastle/pkg/config"
)

// ConfigHandler handles the config requests coming from the frontend
func ConfigHandler(responseWriter http.ResponseWriter, request *http.Request) {
	enableCors(&responseWriter, "*")

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("Failed to read configuration for forecastle", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(appConfig)
	if err != nil {
		logger.Error("An error occurred while marshalling configuration to json", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	_, err = responseWriter.Write(js)
	if err != nil {
		logger.Error("An error occurred while rendering json to output: ", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}

}

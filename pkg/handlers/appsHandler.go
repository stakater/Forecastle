package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stakater/Forecastle/pkg/apps"
	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppsHandler func responsible for serving apps at /apps
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	kubeClient := kube.GetClient()

	var forecastleApps []apps.ForecastleApp
	var err error

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("Failed to read configuration for forecastle", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	appsList := apps.NewList(kubeClient, *appConfig)

	if len(appConfig.Namespaces) != 0 {
		logger.Info("Looking for forecastle apps in these namespaces: ", appConfig.Namespaces)
		forecastleApps, err = appsList.Populate(appConfig.Namespaces...).Get()
	} else {
		logger.Info("Namespaces filter not found. Looking for forecastle apps in all namespaces")
		forecastleApps, err = appsList.Populate(metav1.NamespaceAll).Get()
	}

	if err != nil {
		logger.Error("An error occurred while looking for forceastle apps", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(forecastleApps)
	if err != nil {
		logger.Error("An error occurred while marshalling apps to json", err)
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

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stakater/Forecastle/pkg/forecastle"
	"github.com/stakater/Forecastle/pkg/kube/util"

	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle/customapps"
	"github.com/stakater/Forecastle/pkg/forecastle/ingressapps"
	"github.com/stakater/Forecastle/pkg/kube"
)

// AppsHandler func responsible for serving apps at /apps
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	kubeClient := kube.GetClient()

	var forecastleApps []forecastle.App
	var err error

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("Failed to read configuration for forecastle", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	ingressAppsList := ingressapps.NewList(kubeClient, *appConfig)

	namespaces, err := util.PopulateNamespaceList(appConfig.NamespaceSelector)

	if err != nil {
		logger.Error("An error occurred while populating namespaces", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Namespaces to look for forecastle apps: ", namespaces)
	forecastleApps, err = ingressAppsList.Populate(namespaces...).Get()

	if err != nil {
		logger.Error("An error occurred while looking for forceastle apps", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	customAppsList := customapps.NewList(*appConfig)

	customForecastleApps, err := customAppsList.Populate().Get()
	if err != nil {
		logger.Error("An error occured while populating custom forecastle apps", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}

	// Append both generated and custom apps
	forecastleApps = append(forecastleApps, customForecastleApps...)

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

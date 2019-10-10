package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stakater/Forecastle/pkg/forecastle"
	"github.com/stakater/Forecastle/pkg/kube/util"

	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle/customapps"
	"github.com/stakater/Forecastle/pkg/forecastle/forecastlecrdapps"
	"github.com/stakater/Forecastle/pkg/forecastle/ingressapps"
	"github.com/stakater/Forecastle/pkg/kube"
)

// AppsHandler func responsible for serving apps at /apps
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	kubeClient := kube.GetClient()
	forecastleClient := kube.GetForecastleClient()

	var forecastleApps []forecastle.App
	var err error

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("Failed to read configuration for forecastle", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	namespaces, err := util.PopulateNamespaceList(appConfig.NamespaceSelector)

	if err != nil {
		logger.Error("An error occurred while populating namespaces", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Looking for forecastle apps in the following namespaces: ", namespaces)

	ingressAppsList := ingressapps.NewList(kubeClient, *appConfig)
	forecastleApps, err = ingressAppsList.Populate(namespaces...).Get()

	if err != nil {
		logger.Error("An error occurred while looking for forceastle apps: ", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	customAppsList := customapps.NewList(*appConfig)

	customForecastleApps, err := customAppsList.Populate().Get()
	if err != nil {
		logger.Error("An error occured while populating custom forecastle apps: ", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}

	// Append both generated and custom apps
	forecastleApps = append(forecastleApps, customForecastleApps...)

	forecastleCRDAppsList := forecastlecrdapps.NewList(forecastleClient, *appConfig)
	forecastleCRDApps, err := forecastleCRDAppsList.Populate(namespaces...).Get()

	// Log and proceed with this error
	if err != nil {
		logger.Error("An error occurred while looking for forceastle CRD apps: ", err)
	} else { // Append forecastle CRD apps
		forecastleApps = append(forecastleApps, forecastleCRDApps...)
	}

	js, err := json.Marshal(forecastleApps)
	if err != nil {
		logger.Error("An error occurred while marshalling apps to json: ", err)
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

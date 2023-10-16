package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	"github.com/stakater/Forecastle/v1/pkg/kube/util"

	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/crdapps"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/customapps"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/ingressapps"
	"github.com/stakater/Forecastle/v1/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppsHandler func responsible for serving apps at /apps
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	enableCors(&responseWriter, "*")

	clients := kube.GetClients()

	var forecastleApps []forecastle.App
	var err error

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("Failed to read configuration for forecastle", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	namespaces, err := util.PopulateNamespaceList(clients.KubernetesClient, appConfig.NamespaceSelector)

	if err != nil {
		logger.Error("An error occurred while populating namespaces", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	var namespacesString string
	// All Namespaces are selected
	if len(namespaces) == 1 && namespaces[0] == metav1.NamespaceAll {
		namespacesString = "* (All Namespaces)"
	} else {
		namespacesString = strings.Join(namespaces, ",")
	}

	logger.Info("Looking for forecastle apps in the following namespaces: " + namespacesString)

	ingressAppsList := ingressapps.NewList(clients.KubernetesClient, *appConfig)
	forecastleApps, err = ingressAppsList.Populate(namespaces...).Get()

	if err != nil {
		logger.Error("An error occurred while looking for forecastle apps: ", err)
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

	if appConfig.CRDEnabled {

		forecastleCRDAppsList := crdapps.NewList(clients, *appConfig)
		forecastleCRDApps, err := forecastleCRDAppsList.Populate(namespaces...).Get()

		// Log and proceed with this error
		if err != nil {
			logger.Error("An error occurred while looking for forecastle CRD apps: ", err)
		} else { // Append forecastle CRD apps
			forecastleApps = append(forecastleApps, forecastleCRDApps...)
		}
	}

	if len(forecastleApps) == 0 {
		forecastleApps = []forecastle.App{}
	}
	js, err := json.Marshal(forecastleApps)
	logger.Info(string(js))
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

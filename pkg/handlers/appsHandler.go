package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stakater/Forecastle/pkg/apps"
	"github.com/stakater/Forecastle/pkg/kube"
	"github.com/stakater/Forecastle/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NamespaceSeparator const is used as a separator for namespaces
	NamespaceSeparator = ","
)

var (
	logger = log.New()
)

// AppsHandler func responsible for serving apps at /apps and /apps/{namespaces}
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {

	kubeClient := kube.GetClient()
	appsList := apps.NewList(kubeClient)

	var forecastleApps []apps.ForecastleApp
	var err error

	if namespaces := request.FormValue("namespaces"); namespaces != "" {
		logger.Info("Looking for forecastle apps in these namespaces: ", namespaces)
		forecastleApps, err = appsList.Populate(strings.Split(namespaces, NamespaceSeparator)...).Get()
	} else {
		logger.Info("Namespaces filter not found. Looking for forecastle apps in all namespaces")
		forecastleApps, err = appsList.Populate(metav1.NamespaceAll).Get()
	}
	if err != nil {
		logger.Error("An error occurred while looking for forcastle apps", err)
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

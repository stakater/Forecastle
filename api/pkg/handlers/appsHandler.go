package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stakater/Forecastle/api/pkg/apps"
	"github.com/stakater/Forecastle/api/pkg/kube"
)

const (
	// NamespaceSeparator const is used as a separator for namespaces
	NamespaceSeparator = ","
)

// AppsHandler func responsible for serving apps at /apps and /apps/{namespaces}
func AppsHandler(responseWriter http.ResponseWriter, request *http.Request) {

	kubeClient := kube.GetClient()
	appsLister := apps.NewAppsLister(kubeClient)

	var forecastleApps []apps.ForecastleApp
	var err error

	if namespaces := request.FormValue("namespaces"); namespaces != "" {
		forecastleApps, err = appsLister.List(strings.Split(namespaces, NamespaceSeparator)...)
	} else {
		forecastleApps, err = appsLister.ListAll()
	}
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(forecastleApps)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(js)
}

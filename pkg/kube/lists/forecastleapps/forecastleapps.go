package forecastleapps

import (
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	forecastlev1alpha1 "github.com/stakater/Forecastle/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FilterFunc defined for creating functions that comply with the filtering forecastleapps
type FilterFunc func(v1alpha1.ForecastleApp, config.Config) bool

// List struct is used to list forecastleapps
type List struct {
	appConfig        config.Config
	err              error // Used for forwarding errors
	items            []v1alpha1.ForecastleApp
	forecastleClient forecastlev1alpha1.Interface
}

// NewList creates an List object that you can use to query forecastleapps
func NewList(forecastleClient forecastlev1alpha1.Interface, appConfig config.Config, items ...v1alpha1.ForecastleApp) *List {
	return &List{
		forecastleClient: forecastleClient,
		appConfig:        appConfig,
		items:            items,
	}
}

// Populate function returns a list of forecastleapps
func (il *List) Populate(namespaces ...string) *List {
	for _, namespace := range namespaces {
		forecastleapps, err := il.forecastleClient.ForecastleV1alpha1().ForecastleApps(namespace).List(metav1.ListOptions{})
		if err != nil {
			il.err = err
		}
		il.items = append(il.items, forecastleapps.Items...)
	}

	return il
}

// Filter function applies a filter func that is passed as a parameter to the list of forecastleapps
func (il *List) Filter(filterFunc FilterFunc) *List {
	var filtered []v1alpha1.ForecastleApp

	for _, forecastleApp := range il.items {
		if filterFunc(forecastleApp, il.appConfig) {
			filtered = append(filtered, forecastleApp)
		}
	}

	// Replace original forecastleapps with filtered
	il.items = filtered
	return il
}

// Get function returns the forecastleapps currently present in List
func (il *List) Get() ([]v1alpha1.ForecastleApp, error) {
	return il.items, il.err
}

package httproutes

import (
	"context"

	"github.com/stakater/Forecastle/v1/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gateway "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// List struct is used to list HTTPRoutes
type List struct {
	appConfig     config.Config
	err           error
	items         []gatewayv1.HTTPRoute
	gatewayClient gateway.Interface
}

// FilterFunc defined for creating functions that filter HTTPRoutes
type FilterFunc func(gatewayv1.HTTPRoute, config.Config) bool

// NewList creates a List object to query HTTPRoutes
func NewList(gatewayClient gateway.Interface, appConfig config.Config, items ...gatewayv1.HTTPRoute) *List {
	return &List{
		gatewayClient: gatewayClient,
		appConfig:     appConfig,
		items:         items,
	}
}

// Populate returns a list of HTTPRoutes from the specified namespaces
func (hl *List) Populate(namespaces ...string) *List {
	if hl.gatewayClient == nil {
		return hl
	}

	for _, namespace := range namespaces {
		httpRoutes, err := hl.gatewayClient.GatewayV1().HTTPRoutes(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			hl.err = err
			continue
		}
		hl.items = append(hl.items, httpRoutes.Items...)
	}

	return hl
}

// Filter applies a filter function to the list of HTTPRoutes
func (hl *List) Filter(filterFunc FilterFunc) *List {
	var filtered []gatewayv1.HTTPRoute

	for _, httpRoute := range hl.items {
		if filterFunc(httpRoute, hl.appConfig) {
			filtered = append(filtered, httpRoute)
		}
	}

	hl.items = filtered
	return hl
}

// Get returns the HTTPRoutes currently present in List
func (hl *List) Get() ([]gatewayv1.HTTPRoute, error) {
	return hl.items, hl.err
}

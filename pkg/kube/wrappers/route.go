package wrappers

import (
	routev1 "github.com/openshift/api/route/v1"
)

// RouteWrapper struct wraps an Openshift route object
type RouteWrapper struct {
	route *routev1.Route
}

// NewRouteWrapper func creates an instance of RouteWrapper
func NewRouteWrapper(route *routev1.Route) *RouteWrapper {
	return &RouteWrapper{
		route: route,
	}
}

// GetURL func extracts URL of the route wrapped by the object
func (rw *RouteWrapper) GetURL() string {
	prefix := "http://"
	if rw.route.Spec.TLS != nil {
		prefix = "https://"
	}

	return prefix + rw.route.Spec.Host
}

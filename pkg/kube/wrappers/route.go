package wrappers

import (
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stakater/Forecastle/v1/pkg/annotations"
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

// GetAnnotationValue extracts an annotation's value present on the route wrapped by the object
func (rw *RouteWrapper) GetAnnotationValue(annotationKey string) string {
	return getAnnotationValue(rw.route.Annotations, annotationKey)
}

// GetURL func extracts URL of the route wrapped by the object
func (rw *RouteWrapper) GetURL() string {
	if urlFromAnnotation := rw.GetAnnotationValue(annotations.ForecastleURLAnnotation); urlFromAnnotation != "" {
		return urlFromAnnotation
	}

	prefix := "http://"
	if rw.route.Spec.TLS != nil {
		prefix = "https://"
	}

	return prefix + rw.route.Spec.Host
}

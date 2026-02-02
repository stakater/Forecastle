package wrappers

import (
	"strings"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteWrapper wraps a Gateway API HTTPRoute
type HTTPRouteWrapper struct {
	httpRoute *gatewayv1.HTTPRoute
}

// NewHTTPRouteWrapper creates a new HTTPRouteWrapper
func NewHTTPRouteWrapper(httpRoute *gatewayv1.HTTPRoute) *HTTPRouteWrapper {
	return &HTTPRouteWrapper{httpRoute: httpRoute}
}

// GetAnnotationValue extracts an annotation value from the HTTPRoute
func (hw *HTTPRouteWrapper) GetAnnotationValue(annotationKey string) string {
	return getAnnotationValue(hw.httpRoute.Annotations, annotationKey)
}

// GetName returns the name of the HTTPRoute (from annotation or resource name)
func (hw *HTTPRouteWrapper) GetName() string {
	if nameFromAnnotation := hw.GetAnnotationValue(annotations.ForecastleAppNameAnnotation); nameFromAnnotation != "" {
		return nameFromAnnotation
	}
	return hw.httpRoute.Name
}

// GetNamespace returns the namespace of the HTTPRoute
func (hw *HTTPRouteWrapper) GetNamespace() string {
	return hw.httpRoute.Namespace
}

// GetGroup returns the group name (normalized to lowercase)
func (hw *HTTPRouteWrapper) GetGroup() string {
	if groupFromAnnotation := hw.GetAnnotationValue(annotations.ForecastleGroupAnnotation); groupFromAnnotation != "" {
		return strings.ToLower(groupFromAnnotation)
	}
	return strings.ToLower(hw.GetNamespace())
}

// GetProperties parses custom properties from annotation
func (hw *HTTPRouteWrapper) GetProperties() map[string]string {
	if propertiesFromAnnotation := hw.GetAnnotationValue(annotations.ForecastlePropertiesAnnotation); propertiesFromAnnotation != "" {
		return makeMap(propertiesFromAnnotation)
	}
	return nil
}

// GetURL extracts the URL from the HTTPRoute
func (hw *HTTPRouteWrapper) GetURL() string {
	if urlFromAnnotation := getAndValidateURLAnnotation(hw.httpRoute.Annotations, annotations.ForecastleURLAnnotation); urlFromAnnotation != "" {
		return urlFromAnnotation
	}

	if len(hw.httpRoute.Spec.Hostnames) == 0 {
		logger.Warn("No hostnames defined for HTTPRoute: ", hw.httpRoute.Name)
		return ""
	}

	host := string(hw.httpRoute.Spec.Hostnames[0])
	// TLS is configured on Gateway listener, not HTTPRoute - default to https
	return "https://" + host
}

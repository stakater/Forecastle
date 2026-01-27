package wrappers

import (
	"mvdan.cc/xurls/v2"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	ingressroutev1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
)

// IngressRouteWrapper struct wraps a Traefik ingressroute object
type IngressRouteWrapper struct {
	ingressroute *ingressroutev1.IngressRoute
}

// NewIngressRouteWrapper func creates an instance of IngressRouteWrapper
func NewIngressRouteWrapper(ingressroute *ingressroutev1.IngressRoute) *IngressRouteWrapper {
	return &IngressRouteWrapper{
		ingressroute: ingressroute,
	}
}

// GetAnnotationValue extracts an annotation's value present on the ingressroute wrapped by the object
func (irw *IngressRouteWrapper) GetAnnotationValue(annotationKey string) string {
	return getAnnotationValue(irw.ingressroute.Annotations, annotationKey)
}

// GetURL func extracts URL of the route wrapped by the object
func (irw *IngressRouteWrapper) GetURL() string {
	if urlFromAnnotation := irw.GetAnnotationValue(annotations.ForecastleURLAnnotation); urlFromAnnotation != "" {
		return urlFromAnnotation
	}

	xurlsStrict := xurls.Relaxed()
	parsedUrl := ""

	for _, element := range irw.ingressroute.Spec.Routes {
		tempUrl := xurlsStrict.FindString(element.Match)
		if len(tempUrl) > 0 {
			parsedUrl = tempUrl
		}
	}
	if len(parsedUrl) == 0 {
		logger.Warn("No route url exist in ingressroute: ", irw.ingressroute.GetName())
		return ""
	}

	prefix := "http://"
	if irw.ingressroute.Spec.TLS != nil {
		prefix = "https://"
	}
	return prefix + parsedUrl
}

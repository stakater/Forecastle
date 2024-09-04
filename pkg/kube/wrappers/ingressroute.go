package wrappers

import (
	"mvdan.cc/xurls/v2"

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

// GetURL func extracts URL of the route wrapped by the object
func (irw *IngressRouteWrapper) GetURL() string {
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

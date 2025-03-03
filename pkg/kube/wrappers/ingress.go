package wrappers

import (
	"net/url"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/log"
	v1 "k8s.io/api/networking/v1"
)

var (
	logger = log.New()
)

// IngressWrapper struct wraps a kubernetes ingress object
type IngressWrapper struct {
	ingress *v1.Ingress
}

// NewIngressWrapper func creates an instance of IngressWrapper
func NewIngressWrapper(ingress *v1.Ingress) *IngressWrapper {
	return &IngressWrapper{
		ingress: ingress,
	}
}

// GetAnnotationValue extracts an annotation's value present on the ingress wrapped by the object
func (iw *IngressWrapper) GetAnnotationValue(annotationKey string) string {
	if value, ok := iw.ingress.Annotations[annotationKey]; ok {
		return value
	}
	return ""
}

// GetName func extracts name of the ingress wrapped by the object
func (iw *IngressWrapper) GetName() string {
	if nameFromAnnotation := iw.GetAnnotationValue(annotations.ForecastleAppNameAnnotation); nameFromAnnotation != "" {
		return nameFromAnnotation
	}
	return iw.ingress.ObjectMeta.Name
}

// GetNamespace func extracts namespace of the ingress wrapped by the object
func (iw *IngressWrapper) GetNamespace() string {
	return iw.ingress.ObjectMeta.Namespace
}

// GetGroup func extracts group name from the ingress
func (iw *IngressWrapper) GetGroup() string {
	if groupFromAnnotation := iw.GetAnnotationValue(annotations.ForecastleGroupAnnotation); groupFromAnnotation != "" {
		return groupFromAnnotation
	}
	return iw.GetNamespace()
}

func (iw *IngressWrapper) GetProperties() map[string]string {
	if propertiesFromAnnotations := iw.GetAnnotationValue(annotations.ForecastlePropertiesAnnotation); propertiesFromAnnotations != "" {
		return makeMap(propertiesFromAnnotations)
	}
	return nil
}

// GetURL func extracts url of the ingress wrapped by the object
func (iw *IngressWrapper) GetURL() string {

	if urlFromAnnotation := iw.GetAnnotationValue(annotations.ForecastleURLAnnotation); urlFromAnnotation != "" {
		parsedURL, err := url.ParseRequestURI(urlFromAnnotation)
		if err != nil {
			logger.Warn(err)
			return ""
		}
		return parsedURL.String()
	}

	if !iw.rulesExist() {
		logger.Warn("No rules exist in ingress: ", iw.ingress.GetName())
		return ""
	}

	var url string

	if host, exists := iw.tryGetTLSHost(); exists { // Get TLS Host if it exists
		url = host
	} else {
		url = iw.getHost() // Fallback for normal Host
	}

	// Append port + ingressSubPath
	url += iw.getIngressSubPath()

	return url
}

func (iw *IngressWrapper) rulesExist() bool {
	return len(iw.ingress.Spec.Rules) > 0
}

func (iw *IngressWrapper) tryGetTLSHost() (string, bool) {
	if iw.supportsTLS() && len(iw.ingress.Spec.TLS[0].Hosts) > 0 {
		return "https://" + iw.ingress.Spec.TLS[0].Hosts[0], true
	}

	return "", false
}

func (iw *IngressWrapper) supportsTLS() bool {
	return len(iw.ingress.Spec.TLS) > 0
}

func (iw *IngressWrapper) getHost() string {
	return "http://" + iw.ingress.Spec.Rules[0].Host
}

func (iw *IngressWrapper) getIngressSubPath() string {
	rule := iw.ingress.Spec.Rules[0]
	if rule.HTTP != nil {
		if len(rule.HTTP.Paths) > 0 {
			return rule.HTTP.Paths[0].Path
		}
	}
	return ""
}

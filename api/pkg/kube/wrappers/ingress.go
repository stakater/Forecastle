package wrappers

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/extensions/v1beta1"
)

// IngressWrapper struct wraps a kubernetes ingress object
type IngressWrapper struct {
	ingress *v1beta1.Ingress
}

// NewIngressWrapper func creates an instance of IngressWrapper
func NewIngressWrapper(ingress *v1beta1.Ingress) *IngressWrapper {
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
	return iw.ingress.ObjectMeta.Name
}

// GetNamespace func extracts namespace of the ingress wrapped by the object
func (iw *IngressWrapper) GetNamespace() string {
	return iw.ingress.ObjectMeta.Namespace
}

// GetURL func extracts url of the ingress wrapped by the object
func (iw *IngressWrapper) GetURL() string {

	if !iw.rulesExist() {
		log.Warn("No rules exist in ingress: ", iw.ingress.GetName())
		return ""
	}

	var url string

	if host, exists := iw.tryGetTLSHost(); exists { // Get TLS Host if it exists
		url = host
	} else {
		url = iw.getHost() // Fallback for normal Host
	}

	// Append port + ingressSubPath
	url += iw.getIngressSubPathWithPort()

	return url

}

func (iw *IngressWrapper) rulesExist() bool {
	if iw.ingress.Spec.Rules != nil && len(iw.ingress.Spec.Rules) > 0 {
		return true
	}
	return false
}

func (iw *IngressWrapper) tryGetTLSHost() (string, bool) {
	if iw.supportsTLS() {
		return "https://" + iw.ingress.Spec.TLS[0].Hosts[0], true
	}

	return "", false
}

func (iw *IngressWrapper) supportsTLS() bool {
	if iw.ingress.Spec.TLS != nil && len(iw.ingress.Spec.TLS) > 0 {
		return true
	}
	return false
}

func (iw *IngressWrapper) getHost() string {
	return "http://" + iw.ingress.Spec.Rules[0].Host
}

func (iw *IngressWrapper) getIngressSubPathWithPort() string {
	port := iw.getIngressPort()
	subPath := iw.getIngressSubPath()

	return port + subPath
}

func (iw *IngressWrapper) getIngressPort() string {
	rule := iw.ingress.Spec.Rules[0]
	if rule.HTTP != nil {
		if rule.HTTP.Paths != nil && len(rule.HTTP.Paths) > 0 {
			return rule.HTTP.Paths[0].Backend.ServicePort.StrVal
		}
	}
	return ""
}

func (iw *IngressWrapper) getIngressSubPath() string {
	rule := iw.ingress.Spec.Rules[0]
	if rule.HTTP != nil {
		if rule.HTTP.Paths != nil && len(rule.HTTP.Paths) > 0 {
			return rule.HTTP.Paths[0].Path
		}
	}
	return ""
}

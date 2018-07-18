package apps

import "k8s.io/api/extensions/v1beta1"

// For filtering ingresses having ingress class annotation
func byIngressClassAnnotation(ingress v1beta1.Ingress) bool {
	if _, ok := ingress.Annotations[IngressClassAnnotation]; ok {
		return true
	}
	return false
}

// For filtering ingressing having forecastle expose annotation set to true
func byForecastleExposeAnnotation(ingress v1beta1.Ingress) bool {
	if val, ok := ingress.Annotations[ForecastleExposeAnnotation]; ok {
		// Has Forecastle annotation and is exposed
		if val == "true" {
			return true
		}
	}
	return false
}

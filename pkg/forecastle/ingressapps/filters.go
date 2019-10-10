package ingressapps

import (
	"github.com/stakater/Forecastle/pkg/annotations"
	"github.com/stakater/Forecastle/pkg/util/strings"

	"github.com/stakater/Forecastle/pkg/config"
	"k8s.io/api/extensions/v1beta1"
)

// For filtering ingresses having ingress class annotation
func byIngressClassAnnotation(ingress v1beta1.Ingress, appConfig config.Config) bool {
	if _, ok := ingress.Annotations[annotations.IngressClassAnnotation]; ok {
		return true
	}
	return false
}

// For filtering ingressing having forecastle expose annotation set to true
func byForecastleExposeAnnotation(ingress v1beta1.Ingress, appConfig config.Config) bool {
	if val, ok := ingress.Annotations[annotations.ForecastleExposeAnnotation]; ok {
		// Has Forecastle annotation and is exposed
		if val == "true" {
			return true
		}
	}
	return false
}

// For filtering ingresses by forecastle instance
func byForecastleInstanceAnnotation(ingress v1beta1.Ingress, appConfig config.Config) bool {
	if val, ok := ingress.Annotations[annotations.ForecastleInstanceAnnotation]; ok {
		return strings.ContainsBetweenDelimiter(val, appConfig.InstanceName, ",")
	}
	return false
}

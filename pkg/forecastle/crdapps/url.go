package crdapps

import (
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/kube/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getURL(kubeClient kubernetes.Interface, forecastleApp v1alpha1.ForecastleApp) string {
	if len(forecastleApp.Spec.URL) == 0 {
		return discoverURLFromRefs(kubeClient, forecastleApp)

	}
	return forecastleApp.Spec.URL
}

func discoverURLFromRefs(kubeClient kubernetes.Interface, forecastleApp v1alpha1.ForecastleApp) string {
	urlFrom := forecastleApp.Spec.URLFrom
	if urlFrom == nil {
		logger.Warn("No URL sources set for ForecastleApp: " + forecastleApp.Name)
		return ""
	}

	if urlFrom.IngressRef != nil {
		ingress, err := kubeClient.ExtensionsV1beta1().Ingresses(forecastleApp.Namespace).Get(urlFrom.IngressRef.Name, metav1.GetOptions{})
		if err != nil {
			logger.Warn("Ingress not found with name " + urlFrom.IngressRef.Name)
			return ""
		}
		return wrappers.NewIngressWrapper(ingress).GetURL()
	}

	logger.Warn("Unsupported Ref set on ForecastleApp: " + forecastleApp.Name)
	return ""
}

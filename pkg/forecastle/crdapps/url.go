package crdapps

import (
	routes "github.com/openshift/client-go/route/clientset/versioned"
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/kube"
	"github.com/stakater/Forecastle/pkg/kube/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getURL(clients kube.Clients, forecastleApp v1alpha1.ForecastleApp) string {
	if len(forecastleApp.Spec.URL) == 0 {
		return discoverURLFromRefs(clients, forecastleApp)

	}
	return forecastleApp.Spec.URL
}

func discoverURLFromIngressRef(kubeClient kubernetes.Interface, ingressRef *v1alpha1.IngressURLSource, namespace string) string {
	ingress, err := kubeClient.ExtensionsV1beta1().Ingresses(namespace).Get(ingressRef.Name, metav1.GetOptions{})
	if err != nil {
		logger.Warn("Ingress not found with name " + ingressRef.Name)
		return ""
	}
	return wrappers.NewIngressWrapper(ingress).GetURL()
}

func discoverURLFromRouteRef(routesClient routes.Interface, routeRef *v1alpha1.RouteURLSource, namespace string) string {
	route, err := routesClient.RouteV1().Routes(namespace).Get(routeRef.Name, metav1.GetOptions{})
	if err != nil {
		logger.Warn("Route not found with name " + routeRef.Name)
		return ""
	}

	return wrappers.NewRouteWrapper(route).GetURL()
}

func discoverURLFromRefs(clients kube.Clients, forecastleApp v1alpha1.ForecastleApp) string {
	urlFrom := forecastleApp.Spec.URLFrom
	if urlFrom == nil {
		logger.Warn("No URL sources set for ForecastleApp: " + forecastleApp.Name)
		return ""
	}

	if urlFrom.IngressRef != nil {
		return discoverURLFromIngressRef(clients.KubernetesClient, urlFrom.IngressRef, forecastleApp.Namespace)
	}

	if urlFrom.RouteRef != nil {
		return discoverURLFromRouteRef(clients.RoutesClient, urlFrom.RouteRef, forecastleApp.Namespace)
	}

	logger.Warn("Unsupported Ref set on ForecastleApp: " + forecastleApp.Name)
	return ""
}

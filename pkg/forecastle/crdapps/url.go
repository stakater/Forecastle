package crdapps

import (
	"context"
	"errors"

	routes "github.com/openshift/client-go/route/clientset/versioned"
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/kube"
	"github.com/stakater/Forecastle/pkg/kube/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getURL(clients kube.Clients, forecastleApp v1alpha1.ForecastleApp) (string, error) {
	if len(forecastleApp.Spec.URL) == 0 {
		return discoverURLFromRefs(clients, forecastleApp)

	}
	return forecastleApp.Spec.URL, nil
}

func discoverURLFromIngressRef(kubeClient kubernetes.Interface, ingressRef *v1alpha1.IngressURLSource, namespace string) (string, error) {
	ingress, err := kubeClient.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressRef.Name, metav1.GetOptions{})
	if err != nil {
		logger.Warn("Ingress not found with name " + ingressRef.Name)
		return "", err
	}
	return wrappers.NewIngressWrapper(ingress).GetURL(), nil
}

func discoverURLFromRouteRef(routesClient routes.Interface, routeRef *v1alpha1.RouteURLSource, namespace string) (string, error) {
	route, err := routesClient.RouteV1().Routes(namespace).Get(routeRef.Name, metav1.GetOptions{})
	if err != nil {
		logger.Warn("Route not found with name " + routeRef.Name)
		return "", err
	}

	return wrappers.NewRouteWrapper(route).GetURL(), nil
}

func discoverURLFromRefs(clients kube.Clients, forecastleApp v1alpha1.ForecastleApp) (string, error) {
	urlFrom := forecastleApp.Spec.URLFrom
	if urlFrom == nil {
		logger.Warn("No URL sources set for ForecastleApp: " + forecastleApp.Name)
		return "", errors.New("No URL sources set for ForecastleApp: " + forecastleApp.Name)
	}

	if urlFrom.IngressRef != nil {
		return discoverURLFromIngressRef(clients.KubernetesClient, urlFrom.IngressRef, forecastleApp.Namespace)
	}

	if urlFrom.RouteRef != nil {
		return discoverURLFromRouteRef(clients.RoutesClient, urlFrom.RouteRef, forecastleApp.Namespace)
	}

	logger.Warn("Unsupported Ref set on ForecastleApp: " + forecastleApp.Name)
	return "", errors.New("Unsupported Ref set on ForecastleApp: " + forecastleApp.Name)
}

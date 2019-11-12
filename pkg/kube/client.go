package kube

import (
	"os"

	routesClient "github.com/openshift/client-go/route/clientset/versioned"
	forecastlev1alpha1 "github.com/stakater/Forecastle/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/pkg/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	logger = log.New()
)

// Clients struct exposes interfaces for kubernetes as well as custom resource interfaces
type Clients struct {
	KubernetesClient     kubernetes.Interface
	ForecastleAppsClient forecastlev1alpha1.Interface
	RoutesClient         routesClient.Interface
}

// GetClients returns a `Clients` object containing all available interfaces
func GetClients() Clients {
	return Clients{
		KubernetesClient:     GetKubernetesClient(),
		ForecastleAppsClient: GetForecastleClient(),
		RoutesClient:         GetRoutesClient(),
	}
}

// GetRoutesClient resturns a routes clientset
func GetRoutesClient() routesClient.Interface {
	config := getClientConfig()
	routesClient, err := routesClient.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Can not create routes client: %v", err)
	}

	return routesClient
}

// GetKubernetesClient returns a k8s clientset
func GetKubernetesClient() kubernetes.Interface {
	config := getClientConfig()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Can not create kubernetes client: %v", err)
	}

	return kubeClient
}

// GetForecastleClient returns a forecastle resource clientset
func GetForecastleClient() forecastlev1alpha1.Interface {
	config := getClientConfig()
	forecastleClient, err := forecastlev1alpha1.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Can not create forecastle client: %v", err)
	}

	return forecastleClient
}

func getClientConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		config = getOutOfClusterConfig()
	}

	return config
}

func getOutOfClusterConfig() *rest.Config {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		logger.Fatalf("Cannot get kubernetes config: %v", err)
	}

	return config
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

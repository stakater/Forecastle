package kube

import (
	"os"

	routesClient "github.com/openshift/client-go/route/clientset/versioned"
	forecastlev1alpha1 "github.com/stakater/Forecastle/v1/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/v1/pkg/log"
	ingressroutesClient "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

var (
	logger = log.New()
)

// Clients struct exposes interfaces for kubernetes as well as custom resource interfaces
type Clients struct {
	KubernetesClient     kubernetes.Interface
	ForecastleAppsClient forecastlev1alpha1.Interface
	RoutesClient         routesClient.Interface
	IngressRoutesClient  ingressroutesClient.Interface
	GatewayClient        gatewayClient.Interface
}

// GetClients returns a Clients object with conditionally initialized clients based on API availability
func GetClients() Clients {
	config := getClientConfig()
	availability := DiscoverAPIs(config)

	clients := Clients{
		KubernetesClient:     getKubernetesClient(),
		ForecastleAppsClient: getForecastleClient(),
	}

	if availability.RoutesAvailable {
		logger.Info("OpenShift Route API detected")
		clients.RoutesClient = getRoutesClient(config)
	}

	if availability.IngressRoutesAvailable {
		logger.Info("Traefik IngressRoute API detected")
		clients.IngressRoutesClient = getIngressRoutesClient(config)
	}

	if availability.HTTPRoutesAvailable {
		logger.Info("Gateway API HTTPRoute detected")
		clients.GatewayClient = getGatewayClient(config)
	}

	return clients
}

func getRoutesClient(config *rest.Config) routesClient.Interface {
	client, err := routesClient.NewForConfig(config)
	if err != nil {
		logger.Warnf("Failed to create routes client: %v", err)
		return nil
	}
	return client
}

func getIngressRoutesClient(config *rest.Config) ingressroutesClient.Interface {
	client, err := ingressroutesClient.NewForConfig(config)
	if err != nil {
		logger.Warnf("Failed to create ingressroutes client: %v", err)
		return nil
	}
	return client
}

func getGatewayClient(config *rest.Config) gatewayClient.Interface {
	client, err := gatewayClient.NewForConfig(config)
	if err != nil {
		logger.Warnf("Failed to create gateway client: %v", err)
		return nil
	}
	return client
}

// getKubernetesClient returns a k8s clientset
func getKubernetesClient() kubernetes.Interface {
	config := getClientConfig()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Can not create kubernetes client: %v", err)
	}

	return kubeClient
}

// getForecastleClient returns a forecastle resource clientset
func getForecastleClient() forecastlev1alpha1.Interface {
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

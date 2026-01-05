package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	forecastlev1alpha1 "github.com/stakater/Forecastle/v1/pkg/apis/forecastle/v1alpha1"
	forecastleClient "github.com/stakater/Forecastle/v1/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	testNamespace = "forecastle-e2e"
	apiTimeout    = 10 * time.Second
)

// TestClients holds Kubernetes clients for e2e tests
type TestClients struct {
	Kubernetes  kubernetes.Interface
	Forecastle  forecastleClient.Interface
	Config      *rest.Config
	BaseURL     string
	CleanupFunc func()
}

// SetupClients creates Kubernetes clients from kubeconfig
func SetupClients() (*TestClients, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	fcClient, err := forecastleClient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecastle client: %w", err)
	}

	return &TestClients{
		Kubernetes: kubeClient,
		Forecastle: fcClient,
		Config:     config,
	}, nil
}

func getConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// CreateTestNamespace creates the test namespace
func (tc *TestClients) CreateTestNamespace(ctx context.Context) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "forecastle-e2e",
			},
		},
	}

	_, err := tc.Kubernetes.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create test namespace: %w", err)
	}

	return nil
}

// DeleteTestNamespace removes the test namespace and all resources
func (tc *TestClients) DeleteTestNamespace(ctx context.Context) error {
	err := tc.Kubernetes.CoreV1().Namespaces().Delete(ctx, testNamespace, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete test namespace: %w", err)
	}

	return nil
}

// CreateIngress creates an Ingress resource with forecastle annotations
func (tc *TestClients) CreateIngress(ctx context.Context, name, host string, opts ...IngressOption) (*networkingv1.Ingress, error) {
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
			Annotations: map[string]string{
				annotations.ForecastleExposeAnnotation: "true",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: name,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(ingress)
	}

	created, err := tc.Kubernetes.NetworkingV1().Ingresses(testNamespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create ingress %s: %w", name, err)
	}

	return created, nil
}

// IngressOption is a function that modifies an Ingress
type IngressOption func(*networkingv1.Ingress)

// WithIcon sets the forecastle icon annotation
func WithIcon(icon string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Annotations[annotations.ForecastleIconAnnotation] = icon
	}
}

// WithGroup sets the forecastle group annotation
func WithGroup(group string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Annotations[annotations.ForecastleGroupAnnotation] = group
	}
}

// WithAppName sets the forecastle appName annotation
func WithAppName(name string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Annotations[annotations.ForecastleAppNameAnnotation] = name
	}
}

// WithInstance sets the forecastle instance annotation
func WithInstance(instance string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Annotations[annotations.ForecastleInstanceAnnotation] = instance
	}
}

// WithURL sets the forecastle URL annotation (overrides ingress host)
func WithURL(url string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Annotations[annotations.ForecastleURLAnnotation] = url
	}
}

// WithTLS adds TLS configuration to the ingress
func WithTLS(host string) IngressOption {
	return func(i *networkingv1.Ingress) {
		i.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts: []string{host},
			},
		}
	}
}

// CreateForecastleApp creates a ForecastleApp CRD
func (tc *TestClients) CreateForecastleApp(ctx context.Context, name, url, group, icon string) (*forecastlev1alpha1.ForecastleApp, error) {
	app := &forecastlev1alpha1.ForecastleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: forecastlev1alpha1.ForecastleAppSpec{
			Name:  name,
			URL:   url,
			Group: group,
			Icon:  icon,
		},
	}

	created, err := tc.Forecastle.ForecastleV1alpha1().ForecastleApps(testNamespace).Create(app)
	if err != nil {
		return nil, fmt.Errorf("failed to create ForecastleApp %s: %w", name, err)
	}

	return created, nil
}

// GetApps fetches apps from the Forecastle API
func (tc *TestClients) GetApps(ctx context.Context) ([]forecastle.App, error) {
	url := fmt.Sprintf("%s/api/apps", tc.BaseURL)

	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch apps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var apps []forecastle.App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("failed to decode apps response: %w", err)
	}

	return apps, nil
}

// GetConfig fetches config from the Forecastle API
func (tc *TestClients) GetConfig(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/config", tc.BaseURL)

	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config response: %w", err)
	}

	return config, nil
}

// CheckHealth checks the health endpoint
func (tc *TestClients) CheckHealth(ctx context.Context) error {
	url := fmt.Sprintf("%s/healthz", tc.BaseURL)

	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// CheckReady checks the readiness endpoint
func (tc *TestClients) CheckReady(ctx context.Context) error {
	url := fmt.Sprintf("%s/readyz", tc.BaseURL)

	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check readiness: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("readiness check failed with status %d", resp.StatusCode)
	}

	return nil
}

// WaitForReady waits for the server to become ready
func (tc *TestClients) WaitForReady(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := tc.CheckReady(ctx); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("server did not become ready within %v", timeout)
}

// FindAppByName finds an app by name in the list
func FindAppByName(apps []forecastle.App, name string) *forecastle.App {
	for i := range apps {
		if apps[i].Name == name {
			return &apps[i]
		}
	}
	return nil
}

// FindAppsByGroup finds all apps in a specific group
func FindAppsByGroup(apps []forecastle.App, group string) []forecastle.App {
	var result []forecastle.App
	for _, app := range apps {
		if app.Group == group {
			result = append(result, app)
		}
	}
	return result
}

package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	v1alpha1 "github.com/stakater/Forecastle/v1/pkg/apis/forecastle/v1alpha1"
	forecastlefake "github.com/stakater/Forecastle/v1/pkg/client/clientset/versioned/fake"
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	"github.com/stakater/Forecastle/v1/pkg/kube"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	gatewayfake "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/fake"
)

func TestHandler_DiscoverApps_IngressOnly(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create an ingress with forecastle annotations
	ingress := testutil.AddAnnotationToIngress(
		testutil.AddAnnotationToIngress(
			testutil.AddAnnotationToIngress(
				testutil.CreateIngressWithHost("my-app", "myapp.example.com"),
				annotations.ForecastleExposeAnnotation, "true"),
			annotations.ForecastleIconAnnotation, "https://example.com/icon.png"),
		annotations.ForecastleGroupAnnotation, "Production")
	ingress.Namespace = "default"
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil, // No Gateway API
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}

	app := apps[0]
	if app.Name != "my-app" {
		t.Errorf("Expected app name 'my-app', got '%s'", app.Name)
	}
	if app.URL != "http://myapp.example.com" {
		t.Errorf("Expected URL 'http://myapp.example.com', got '%s'", app.URL)
	}
	if app.Group != "production" {
		t.Errorf("Expected group 'production', got '%s'", app.Group)
	}
	if app.Icon != "https://example.com/icon.png" {
		t.Errorf("Expected icon 'https://example.com/icon.png', got '%s'", app.Icon)
	}
	if app.DiscoverySource != forecastle.Ingress {
		t.Errorf("Expected discovery source Ingress, got %v", app.DiscoverySource)
	}
}

func TestHandler_DiscoverApps_HTTPRouteOnly(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()
	gatewayClient := gatewayfake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create an HTTPRoute with forecastle annotations
	httpRoute := testutil.AddAnnotationToHTTPRoute(
		testutil.AddAnnotationToHTTPRoute(
			testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRouteWithHostnameAndNamespace("gateway-app", "default", "gateway.example.com"),
				annotations.ForecastleExposeAnnotation, "true"),
			annotations.ForecastleIconAnnotation, "https://example.com/gateway-icon.png"),
		annotations.ForecastleGroupAnnotation, "Gateway")
	_, _ = gatewayClient.GatewayV1().HTTPRoutes("default").Create(context.TODO(), httpRoute, metav1.CreateOptions{})

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        gatewayClient,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}

	app := apps[0]
	if app.Name != "gateway-app" {
		t.Errorf("Expected app name 'gateway-app', got '%s'", app.Name)
	}
	if app.URL != "https://gateway.example.com" {
		t.Errorf("Expected URL 'https://gateway.example.com', got '%s'", app.URL)
	}
	if app.Group != "gateway" {
		t.Errorf("Expected group 'gateway', got '%s'", app.Group)
	}
	if app.DiscoverySource != forecastle.HTTPRoute {
		t.Errorf("Expected discovery source HTTPRoute, got %v", app.DiscoverySource)
	}
}

func TestHandler_DiscoverApps_CRDOnly(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create a ForecastleApp CRD
	forecastleApp := &v1alpha1.ForecastleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "crd-app",
			Namespace: "default",
		},
		Spec: v1alpha1.ForecastleAppSpec{
			Name:  "CRD Application",
			URL:   "https://crd.example.com",
			Icon:  "https://example.com/crd-icon.png",
			Group: "CRD Apps",
		},
	}
	_, _ = forecastleClient.ForecastleV1alpha1().ForecastleApps("default").Create(forecastleApp)

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
		CRDEnabled:        true,
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}

	app := apps[0]
	if app.Name != "CRD Application" {
		t.Errorf("Expected app name 'CRD Application', got '%s'", app.Name)
	}
	if app.URL != "https://crd.example.com" {
		t.Errorf("Expected URL 'https://crd.example.com', got '%s'", app.URL)
	}
	if app.Group != "crd apps" {
		t.Errorf("Expected group 'crd apps', got '%s'", app.Group)
	}
	if app.DiscoverySource != forecastle.ForecastleAppCRD {
		t.Errorf("Expected discovery source ForecastleAppCRD, got %v", app.DiscoverySource)
	}
}

func TestHandler_DiscoverApps_CustomAppsOnly(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
		CustomApps: []config.CustomApp{
			{
				Name:  "Custom App",
				URL:   "https://custom.example.com",
				Icon:  "https://example.com/custom-icon.png",
				Group: "Custom",
			},
		},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}

	app := apps[0]
	if app.Name != "Custom App" {
		t.Errorf("Expected app name 'Custom App', got '%s'", app.Name)
	}
	if app.URL != "https://custom.example.com" {
		t.Errorf("Expected URL 'https://custom.example.com', got '%s'", app.URL)
	}
	if app.Group != "custom" {
		t.Errorf("Expected group 'custom', got '%s'", app.Group)
	}
	if app.DiscoverySource != forecastle.Config {
		t.Errorf("Expected discovery source Config, got %v", app.DiscoverySource)
	}
}

func TestHandler_DiscoverApps_AllSources(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()
	gatewayClient := gatewayfake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create Ingress app
	ingress := testutil.AddAnnotationToIngress(
		testutil.CreateIngressWithHost("ingress-app", "ingress.example.com"),
		annotations.ForecastleExposeAnnotation, "true")
	ingress.Namespace = "default"
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})

	// Create HTTPRoute app
	httpRoute := testutil.AddAnnotationToHTTPRoute(
		testutil.CreateHTTPRouteWithHostnameAndNamespace("httproute-app", "default", "httproute.example.com"),
		annotations.ForecastleExposeAnnotation, "true")
	_, _ = gatewayClient.GatewayV1().HTTPRoutes("default").Create(context.TODO(), httpRoute, metav1.CreateOptions{})

	// Create ForecastleApp CRD
	forecastleApp := &v1alpha1.ForecastleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "crd-app",
			Namespace: "default",
		},
		Spec: v1alpha1.ForecastleAppSpec{
			Name:  "CRD App",
			URL:   "https://crd.example.com",
			Group: "CRD",
		},
	}
	_, _ = forecastleClient.ForecastleV1alpha1().ForecastleApps("default").Create(forecastleApp)

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        gatewayClient,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
		CRDEnabled:        true,
		CustomApps: []config.CustomApp{
			{
				Name:  "Custom App",
				URL:   "https://custom.example.com",
				Group: "Custom",
			},
		},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 4 {
		t.Fatalf("Expected 4 apps (1 ingress + 1 httproute + 1 crd + 1 custom), got %d", len(apps))
	}

	// Verify we have one app from each source
	sources := make(map[forecastle.DiscoverySource]int)
	for _, app := range apps {
		sources[app.DiscoverySource]++
	}

	if sources[forecastle.Ingress] != 1 {
		t.Errorf("Expected 1 Ingress app, got %d", sources[forecastle.Ingress])
	}
	if sources[forecastle.HTTPRoute] != 1 {
		t.Errorf("Expected 1 HTTPRoute app, got %d", sources[forecastle.HTTPRoute])
	}
	if sources[forecastle.ForecastleAppCRD] != 1 {
		t.Errorf("Expected 1 ForecastleAppCRD app, got %d", sources[forecastle.ForecastleAppCRD])
	}
	if sources[forecastle.Config] != 1 {
		t.Errorf("Expected 1 Config app, got %d", sources[forecastle.Config])
	}
}

func TestHandler_DiscoverApps_NamespaceFiltering(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()
	gatewayClient := gatewayfake.NewSimpleClientset()

	// Create namespaces
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "production"}},
		metav1.CreateOptions{},
	)
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "staging"}},
		metav1.CreateOptions{},
	)

	// Create ingress in production
	ingressProd := testutil.AddAnnotationToIngress(
		testutil.CreateIngressWithHost("prod-app", "prod.example.com"),
		annotations.ForecastleExposeAnnotation, "true")
	ingressProd.Namespace = "production"
	_, _ = kubeClient.NetworkingV1().Ingresses("production").Create(context.TODO(), ingressProd, metav1.CreateOptions{})

	// Create ingress in staging
	ingressStaging := testutil.AddAnnotationToIngress(
		testutil.CreateIngressWithHost("staging-app", "staging.example.com"),
		annotations.ForecastleExposeAnnotation, "true")
	ingressStaging.Namespace = "staging"
	_, _ = kubeClient.NetworkingV1().Ingresses("staging").Create(context.TODO(), ingressStaging, metav1.CreateOptions{})

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        gatewayClient,
	}

	// Only select production namespace
	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{
			MatchNames: []string{"production"},
		},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app (only production), got %d", len(apps))
	}

	if apps[0].Name != "prod-app" {
		t.Errorf("Expected 'prod-app', got '%s'", apps[0].Name)
	}
}

func TestHandler_DiscoverApps_InstanceFiltering(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()         //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create ingress for instance "forecastle-a"
	ingressA := testutil.AddAnnotationToIngress(
		testutil.AddAnnotationToIngress(
			testutil.CreateIngressWithHost("app-a", "a.example.com"),
			annotations.ForecastleExposeAnnotation, "true"),
		annotations.ForecastleInstanceAnnotation, "forecastle-a")
	ingressA.Namespace = "default"
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingressA, metav1.CreateOptions{})

	// Create ingress for instance "forecastle-b"
	ingressB := testutil.AddAnnotationToIngress(
		testutil.AddAnnotationToIngress(
			testutil.CreateIngressWithHost("app-b", "b.example.com"),
			annotations.ForecastleExposeAnnotation, "true"),
		annotations.ForecastleInstanceAnnotation, "forecastle-b")
	ingressB.Namespace = "default"
	ingressB.Name = "app-b"
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingressB, metav1.CreateOptions{})

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil,
	}

	// Only select instance "forecastle-a"
	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
		InstanceName:      "forecastle-a",
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	if err != nil {
		t.Fatalf("discoverApps() error = %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app (only forecastle-a instance), got %d", len(apps))
	}

	if apps[0].Name != "app-a" {
		t.Errorf("Expected 'app-a', got '%s'", apps[0].Name)
	}
}

func TestHandler_DiscoverApps_WithoutGatewayClient(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil, // Gateway API not available
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	apps, err := handler.discoverApps(cfg)

	// Should not error even without Gateway client
	if err != nil {
		t.Fatalf("discoverApps() should not error without Gateway client, got %v", err)
	}

	if apps == nil {
		t.Error("apps should not be nil")
	}
}

func TestHandler_AppsHandler_ReturnsJSON(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	// Create namespace
	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	// Create an ingress
	ingress := testutil.AddAnnotationToIngress(
		testutil.CreateIngressWithHost("test-app", "test.example.com"),
		annotations.ForecastleExposeAnnotation, "true")
	ingress.Namespace = "default"
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
		GatewayClient:        nil,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)

	// Manually populate the cache
	handler.refreshCache(context.Background())

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
	rec := httptest.NewRecorder()

	handler.AppsHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}

	// Parse response as generic JSON to verify structure
	var apps []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&apps); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(apps))
	}

	// Verify expected fields
	if apps[0]["name"] != "test-app" {
		t.Errorf("Expected app name 'test-app', got '%v'", apps[0]["name"])
	}
	if apps[0]["url"] != "http://test.example.com" {
		t.Errorf("Expected URL 'http://test.example.com', got '%v'", apps[0]["url"])
	}
	if apps[0]["discoverySource"] != "Ingress" {
		t.Errorf("Expected discoverySource 'Ingress', got '%v'", apps[0]["discoverySource"])
	}
}

func TestHandler_HealthzHandler(t *testing.T) {
	handler := NewHandler(nil, nil, time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.HealthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", rec.Body.String())
	}
}

func TestHandler_ReadyzHandler_NotReady(t *testing.T) {
	handler := NewHandler(nil, nil, time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ReadyzHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", rec.Code)
	}
}

func TestHandler_ReadyzHandler_Ready(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() //nolint:staticcheck // NewClientset requires generated apply configurations
	forecastleClient := forecastlefake.NewSimpleClientset()

	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
		metav1.CreateOptions{},
	)

	clients := &kube.Clients{
		KubernetesClient:     kubeClient,
		ForecastleAppsClient: forecastleClient,
	}

	cfg := &config.Config{
		NamespaceSelector: config.NamespaceSelector{Any: true},
	}

	handler := NewHandler(clients, func() (*config.Config, error) { return cfg, nil }, time.Minute)
	handler.refreshCache(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ReadyzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

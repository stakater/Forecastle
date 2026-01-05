// Package e2e contains end-to-end tests for Forecastle.
// These tests run against a real Kubernetes cluster (typically Kind) and verify
// that Forecastle correctly discovers apps from various sources.
//
// Prerequisites:
//   - A running Kubernetes cluster (Kind recommended)
//   - KUBECONFIG set or ~/.kube/config available
//   - FORECASTLE_URL environment variable pointing to a running Forecastle instance
//     OR run `make build && ./forecastle` before running tests
//
// Run with: go test -v -tags=e2e ./e2e/...
package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testClients *TestClients

var clusterAvailable bool

func TestMain(m *testing.M) {
	// Get Forecastle URL from environment or use default
	baseURL := os.Getenv("FORECASTLE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	clients, err := SetupClients()
	if err != nil {
		// Cluster not available, but we can still run API-only tests
		clients = &TestClients{}
		clusterAvailable = false
	} else {
		clusterAvailable = true
		// Create test namespace once for all tests
		ctx := context.Background()
		if err := clients.CreateTestNamespace(ctx); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				// Log but don't fail - tests will skip if needed
				clusterAvailable = false
			}
		}
	}
	clients.BaseURL = baseURL
	testClients = clients

	// Run tests
	code := m.Run()

	// Cleanup namespace after all tests
	if clusterAvailable && testClients != nil {
		_ = testClients.DeleteTestNamespace(context.Background())
	}

	os.Exit(code)
}

func skipIfNoCluster(t *testing.T) {
	if !clusterAvailable {
		t.Skip("Kubernetes cluster not available, skipping test")
	}
}

// normalizeURL removes trailing slashes for comparison
func normalizeURL(url string) string {
	return strings.TrimSuffix(url, "/")
}

func TestHealthEndpoint(t *testing.T) {
	ctx := context.Background()

	if err := testClients.CheckHealth(ctx); err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
}

func TestReadinessEndpoint(t *testing.T) {
	ctx := context.Background()

	// Wait for server to be ready (cache populated)
	if err := testClients.WaitForReady(ctx, 30*time.Second); err != nil {
		t.Fatalf("Readiness check failed: %v", err)
	}
}

func TestConfigEndpoint(t *testing.T) {
	ctx := context.Background()

	config, err := testClients.GetConfig(ctx)
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	// Verify config has expected structure
	if _, ok := config["title"]; !ok {
		t.Error("Config should have 'title' field")
	}
}

func TestAppsEndpoint(t *testing.T) {
	ctx := context.Background()

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	// Apps should be a valid array (may be empty)
	if apps == nil {
		t.Error("Apps should not be nil")
	}
}

func TestIngressDiscovery(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	// Create test ingress
	ingressName := "e2e-test-app"
	ingressHost := "e2e-test.example.com"
	_, err := testClients.CreateIngress(ctx, ingressName, ingressHost,
		WithGroup("e2e-tests"),
		WithIcon("https://example.com/icon.png"),
	)
	if err != nil {
		t.Fatalf("Failed to create test ingress: %v", err)
	}

	// Wait for cache to refresh (default is 20 seconds)
	t.Log("Waiting for cache refresh...")
	time.Sleep(25 * time.Second)

	// Fetch apps and verify our ingress is discovered
	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	app := FindAppByName(apps, ingressName)
	if app == nil {
		t.Fatalf("Expected to find app '%s' in discovered apps", ingressName)
	}

	// Verify app properties (normalize URLs to handle trailing slash differences)
	expectedURL := "http://" + ingressHost
	if normalizeURL(app.URL) != normalizeURL(expectedURL) {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, app.URL)
	}

	if app.Group != "e2e-tests" {
		t.Errorf("Expected group 'e2e-tests', got '%s'", app.Group)
	}

	if app.Icon != "https://example.com/icon.png" {
		t.Errorf("Expected icon 'https://example.com/icon.png', got '%s'", app.Icon)
	}

	if app.DiscoverySource != forecastle.Ingress {
		t.Errorf("Expected discovery source 'Ingress', got '%v'", app.DiscoverySource)
	}
}

func TestIngressWithTLS(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	// Create TLS ingress
	ingressName := "e2e-tls-app"
	ingressHost := "e2e-tls.example.com"
	_, err := testClients.CreateIngress(ctx, ingressName, ingressHost,
		WithTLS(ingressHost),
		WithGroup("e2e-tls"),
	)
	if err != nil {
		t.Fatalf("Failed to create TLS ingress: %v", err)
	}

	// Wait for cache refresh
	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	app := FindAppByName(apps, ingressName)
	if app == nil {
		t.Fatalf("Expected to find app '%s'", ingressName)
	}

	// TLS ingresses should have https URL
	if !strings.HasPrefix(app.URL, "https://") {
		t.Errorf("Expected HTTPS URL for TLS ingress, got '%s'", app.URL)
	}
}

func TestIngressWithCustomAppName(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	ingressName := "e2e-custom-name-ingress"
	customAppName := "My Custom App Name"
	_, err := testClients.CreateIngress(ctx, ingressName, "custom.example.com",
		WithAppName(customAppName),
	)
	if err != nil {
		t.Fatalf("Failed to create ingress: %v", err)
	}

	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	// Should be found by custom app name, not ingress name
	app := FindAppByName(apps, customAppName)
	if app == nil {
		t.Fatalf("Expected to find app with custom name '%s'", customAppName)
	}
}

func TestIngressWithURLOverride(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	ingressName := "e2e-url-override"
	customURL := "https://custom-url.example.com/app"
	_, err := testClients.CreateIngress(ctx, ingressName, "ingress-host.example.com",
		WithURL(customURL),
	)
	if err != nil {
		t.Fatalf("Failed to create ingress: %v", err)
	}

	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	app := FindAppByName(apps, ingressName)
	if app == nil {
		t.Fatalf("Expected to find app '%s'", ingressName)
	}

	// URL should be the override, not the ingress host
	if app.URL != customURL {
		t.Errorf("Expected URL '%s', got '%s'", customURL, app.URL)
	}
}

func TestForecastleAppCRD(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	// Skip if CRD is not installed
	_, err := testClients.Forecastle.ForecastleV1alpha1().ForecastleApps(testNamespace).List(metav1.ListOptions{})
	if err != nil {
		t.Skipf("ForecastleApp CRD not available, skipping: %v", err)
	}

	// Create ForecastleApp CRD
	appName := "e2e-crd-app"
	_, err = testClients.CreateForecastleApp(ctx, appName,
		"https://crd-app.example.com",
		"e2e-crd-group",
		"https://example.com/crd-icon.png",
	)
	if err != nil {
		t.Fatalf("Failed to create ForecastleApp: %v", err)
	}

	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	app := FindAppByName(apps, appName)
	if app == nil {
		t.Fatalf("Expected to find CRD app '%s'", appName)
	}

	if app.DiscoverySource != forecastle.ForecastleAppCRD {
		t.Errorf("Expected discovery source 'ForecastleAppCRD', got '%v'", app.DiscoverySource)
	}

	if app.URL != "https://crd-app.example.com" {
		t.Errorf("Expected URL 'https://crd-app.example.com', got '%s'", app.URL)
	}
}

func TestMultipleAppsInSameGroup(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	groupName := "e2e-multi-app-group"

	// Create multiple apps in the same group
	for i := 1; i <= 3; i++ {
		name := "e2e-multi-app-" + string(rune('0'+i))
		host := name + ".example.com"
		_, err := testClients.CreateIngress(ctx, name, host,
			WithGroup(groupName),
		)
		if err != nil {
			t.Fatalf("Failed to create ingress %s: %v", name, err)
		}
	}

	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	groupApps := FindAppsByGroup(apps, groupName)
	if len(groupApps) != 3 {
		t.Errorf("Expected 3 apps in group '%s', got %d", groupName, len(groupApps))
	}
}

func TestAppWithoutExposeAnnotation(t *testing.T) {
	skipIfNoCluster(t)
	ctx := context.Background()

	// Create ingress WITHOUT the expose annotation
	ingressName := "e2e-not-exposed"
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: testNamespace,
			// No forecastle annotations!
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "not-exposed.example.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: ingressName,
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

	_, err := testClients.Kubernetes.NetworkingV1().Ingresses(testNamespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create ingress: %v", err)
	}

	time.Sleep(25 * time.Second)

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	// This app should NOT be discovered
	app := FindAppByName(apps, ingressName)
	if app != nil {
		t.Errorf("Ingress without expose annotation should not be discovered")
	}
}

func TestCustomAppsFromConfig(t *testing.T) {
	// Custom apps are defined in config.yaml and don't require a cluster
	ctx := context.Background()

	// Wait for cache to be populated
	if err := testClients.WaitForReady(ctx, 30*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	// Find apps with Config discovery source
	var configApps []forecastle.App
	for _, app := range apps {
		if app.DiscoverySource == forecastle.Config {
			configApps = append(configApps, app)
		}
	}

	// Check that at least one custom app from config exists
	// The config.yaml should have custom apps defined
	if len(configApps) == 0 {
		t.Log("No custom apps found with Config discovery source - this is expected if config.yaml has no customApps")
		return
	}

	t.Logf("Found %d custom apps from config", len(configApps))

	// Verify custom apps have required fields
	for _, app := range configApps {
		if app.Name == "" {
			t.Errorf("Custom app should have a name")
		}
		if app.URL == "" {
			t.Errorf("Custom app '%s' should have a URL", app.Name)
		}
		if app.DiscoverySource != forecastle.Config {
			t.Errorf("Custom app '%s' should have discovery source 'Config', got '%v'", app.Name, app.DiscoverySource)
		}
	}
}

func TestAllDiscoverySourcesHaveValidFormat(t *testing.T) {
	// This test verifies that all apps have valid discovery sources
	ctx := context.Background()

	if err := testClients.WaitForReady(ctx, 30*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	apps, err := testClients.GetApps(ctx)
	if err != nil {
		t.Fatalf("Failed to get apps: %v", err)
	}

	// Count apps by discovery source
	sourceCounts := make(map[forecastle.DiscoverySource]int)
	for _, app := range apps {
		sourceCounts[app.DiscoverySource]++

		// Verify each app has required fields regardless of source
		if app.Name == "" {
			t.Errorf("App should have a name (source: %v)", app.DiscoverySource)
		}
		if app.URL == "" {
			t.Errorf("App '%s' should have a URL (source: %v)", app.Name, app.DiscoverySource)
		}
	}

	// Log discovery source distribution
	t.Logf("Discovery source distribution:")
	for source, count := range sourceCounts {
		t.Logf("  %v: %d apps", source, count)
	}
}

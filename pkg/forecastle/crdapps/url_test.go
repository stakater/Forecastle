package crdapps

import (
	"testing"

	routefake "github.com/openshift/client-go/route/clientset/versioned/fake"
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/kube"
	"github.com/stakater/Forecastle/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

func Test_getURL(t *testing.T) {
	clients := kube.Clients{
		RoutesClient:     routefake.NewSimpleClientset(),
		KubernetesClient: kubefake.NewSimpleClientset(),
	}
	type args struct {
		clients       kube.Clients
		forecastleApp v1alpha1.ForecastleApp
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetURLWithDefaultURLValue",
			args: args{
				clients:       clients,
				forecastleApp: *testutil.CreateForecastleApp("app-1", "https://google.com", "default", "https://icon"),
			},
			want: "https://google.com",
		},
		{
			name: "TestGetURLWithNoURL",
			args: args{
				clients:       clients,
				forecastleApp: *testutil.CreateForecastleApp("app-1", "", "default", "https://icon"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getURL(tt.args.clients, tt.args.forecastleApp); got != tt.want {
				t.Errorf("getURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverURLFromRefs(t *testing.T) {
	clients := kube.Clients{
		RoutesClient:     routefake.NewSimpleClientset(),
		KubernetesClient: kubefake.NewSimpleClientset(),
	}
	clients.KubernetesClient.ExtensionsV1beta1().Ingresses("").Create(testutil.CreateIngressWithHost("my-app-ingress", "https://ingress-url.com"))
	clients.RoutesClient.RouteV1().Routes("").Create(testutil.CreateRouteWithHost("my-app-route", "ingress-url.com"))
	type args struct {
		clients       kube.Clients
		forecastleApp v1alpha1.ForecastleApp
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDiscoverURLFromRefsWithIngressName",
			args: args{
				clients:       clients,
				forecastleApp: *testutil.CreateForecastleAppWithURLFromIngress("app-1", "default", "https://icon", "my-app-ingress"),
			},
			want: "http://https://ingress-url.com",
		},
		{
			name: "TestDiscoverURLFromRefsWithRouteName",
			args: args{
				clients:       clients,
				forecastleApp: *testutil.CreateForecastleAppWithURLFromRoute("app-1", "default", "https://icon", "my-app-route"),
			},
			want: "http://ingress-url.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discoverURLFromRefs(tt.args.clients, tt.args.forecastleApp); got != tt.want {
				t.Errorf("discoverURLFromRefs() = %v, want %v", got, tt.want)
			}
		})
	}

	clients.KubernetesClient.ExtensionsV1beta1().Ingresses("").Delete("my-app-ingress", &metav1.DeleteOptions{})
	clients.RoutesClient.RouteV1().Routes("").Delete("my-app-route", &metav1.DeleteOptions{})
}

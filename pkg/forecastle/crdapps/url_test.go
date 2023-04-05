package crdapps

import (
	"context"
	"testing"

	routefake "github.com/openshift/client-go/route/clientset/versioned/fake"
	v1alpha1 "github.com/stakater/Forecastle/v1/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/v1/pkg/kube"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
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
		err  error
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
			if got, err := getURL(tt.args.clients, tt.args.forecastleApp); got != tt.want && err != tt.err {
				t.Errorf("getURL() = %v, want %v, err = %v, wantErr = %v", got, tt.want, err, tt.err)
			}
		})
	}
}

func Test_discoverURLFromRefs(t *testing.T) {
	clients := kube.Clients{
		RoutesClient:     routefake.NewSimpleClientset(),
		KubernetesClient: kubefake.NewSimpleClientset(),
	}
	_, err := clients.KubernetesClient.NetworkingV1().Ingresses("").Create(context.TODO(), testutil.CreateIngressWithHost("my-app-ingress", "https://ingress-url.com"), metav1.CreateOptions{})
	if err != nil {
		t.Errorf("CreateIngressWithHost(\"my-app-ingress\", \"https://ingress-url.com\") Failed")
	}
	_, err = clients.RoutesClient.RouteV1().Routes("").Create(context.TODO(), testutil.CreateRouteWithHost("my-app-route", "ingress-url.com"), metav1.CreateOptions{})
	if err != nil {
		t.Errorf("CreateRouteWithHost(\"my-app-route\", \"ingress-url.com\") Failed")
	}
	type args struct {
		clients       kube.Clients
		forecastleApp v1alpha1.ForecastleApp
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
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
			if got, err := discoverURLFromRefs(tt.args.clients, tt.args.forecastleApp); got != tt.want && err != tt.err {
				t.Errorf("discoverURLFromRefs() = %v, want %v, err = %v, wantErr = %v", got, tt.want, err, tt.err)
			}
		})
	}

	err = clients.KubernetesClient.NetworkingV1().Ingresses("").Delete(context.TODO(), "my-app-ingress", metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("Deleting my-app-ingress Failed")
	}
	err = clients.RoutesClient.RouteV1().Routes("").Delete(context.TODO(), "my-app-route", metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("Deleting my-app-route Failed")
	}
}

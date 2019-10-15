package crdapps

import (
	"testing"

	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_getURL(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	type args struct {
		kubeClient    kubernetes.Interface
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
				kubeClient:    kubeClient,
				forecastleApp: *testutil.CreateForecastleApp("app-1", "https://google.com", "default", "https://icon"),
			},
			want: "https://google.com",
		},
		{
			name: "TestGetURLWithNoURL",
			args: args{
				kubeClient:    kubeClient,
				forecastleApp: *testutil.CreateForecastleApp("app-1", "", "default", "https://icon"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getURL(tt.args.kubeClient, tt.args.forecastleApp); got != tt.want {
				t.Errorf("getURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverURLFromRefs(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	kubeClient.ExtensionsV1beta1().Ingresses("").Create(testutil.CreateIngressWithHost("my-app-ingress", "https://ingress-url.com"))

	type args struct {
		kubeClient    kubernetes.Interface
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
				kubeClient:    kubeClient,
				forecastleApp: *testutil.CreateForecastleAppWithURLFromIngress("app-1", "default", "https://icon", "my-app-ingress"),
			},
			want: "http://https://ingress-url.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discoverURLFromRefs(tt.args.kubeClient, tt.args.forecastleApp); got != tt.want {
				t.Errorf("discoverURLFromRefs() = %v, want %v", got, tt.want)
			}
		})
	}

	kubeClient.ExtensionsV1beta1().Ingresses("").Delete("my-app-ingress", &metav1.DeleteOptions{})

}

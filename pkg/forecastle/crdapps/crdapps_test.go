package crdapps

import (
	"reflect"
	"testing"

	routefake "github.com/openshift/client-go/route/clientset/versioned/fake"
	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/client/clientset/versioned/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"

	"github.com/stakater/Forecastle/pkg/kube"

	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle"
	"github.com/stakater/Forecastle/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewList(t *testing.T) {
	clients := kube.Clients{
		ForecastleAppsClient: fake.NewSimpleClientset(),
		KubernetesClient:     kubefake.NewSimpleClientset(),
	}

	type args struct {
		clients   kube.Clients
		appConfig config.Config
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithNokubeClient",
			args: args{},
			want: &List{},
		},
		{
			name: "TestNewListWithDefaultkubeClient",
			args: args{
				clients: clients,
			},
			want: &List{
				clients: clients,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.clients, tt.args.appConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Populate(t *testing.T) {

	clients := kube.Clients{
		ForecastleAppsClient: fake.NewSimpleClientset(),
		KubernetesClient:     kubefake.NewSimpleClientset(),
	}

	forecastleApp := testutil.CreateForecastleApp("app-1", "https://google.com", "default", "https://google.com/icon.png")

	_, _ = clients.ForecastleAppsClient.ForecastleV1alpha1().ForecastleApps("default").Create(forecastleApp)

	type args struct {
		namespaces []string
	}
	tests := []struct {
		name string
		al   *List
		args args
		want *List
	}{
		{
			name: "TestListPopulateWithSelectedNamespaces",
			al: &List{
				clients: clients,
			},
			args: args{
				namespaces: []string{"default"},
			},
			want: &List{
				clients: clients,
				items: []forecastle.App{
					{
						Name:            "app-1",
						Group:           "default",
						URL:             "https://google.com",
						Icon:            "https://google.com/icon.png",
						DiscoverySource: forecastle.ForecastleAppCRD,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.al.Populate(tt.args.namespaces...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Populate() = %v, want %v", got, tt.want)
			}
		})
	}

	_ = clients.ForecastleAppsClient.ForecastleV1alpha1().ForecastleApps("default").Delete("app-1", &metav1.DeleteOptions{})

}

func TestList_Get(t *testing.T) {
	clients := kube.Clients{
		ForecastleAppsClient: fake.NewSimpleClientset(),
		KubernetesClient:     kubefake.NewSimpleClientset(),
	}

	tests := []struct {
		name    string
		al      *List
		want    []forecastle.App
		wantErr bool
	}{
		{
			name: "TestGetForecastleApps",
			al: &List{
				clients: clients,
				items: []forecastle.App{
					{
						Name:  "app",
						Icon:  "https://google.com/icon.png",
						Group: "test",
						URL:   "https://google.com",
					},
				},
			},
			want: []forecastle.App{
				{
					Name:  "app",
					Icon:  "https://google.com/icon.png",
					Group: "test",
					URL:   "https://google.com",
				},
			},
			wantErr: false,
		},
		{
			name: "TestGetForecastleAppsWithEmptyList",
			al: &List{
				clients: clients,
				items:   []forecastle.App{},
			},
			want:    []forecastle.App{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.al.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertForecastleAppCustomResourcesToForecastleApps(t *testing.T) {
	clients := kube.Clients{
		ForecastleAppsClient: fake.NewSimpleClientset(),
		KubernetesClient:     kubefake.NewSimpleClientset(),
		RoutesClient:         routefake.NewSimpleClientset(),
	}

	type args struct {
		forecastleApps []v1alpha1.ForecastleApp
	}
	tests := []struct {
		name     string
		args     args
		wantApps []forecastle.App
		err      error
	}{
		{
			name: "TestConvertForecastleAppCustomResourcesToForecastleAppsWithNoApps",
			args: args{
				forecastleApps: []v1alpha1.ForecastleApp{},
			},
			wantApps: nil,
		},
		{
			name: "TestConvertForecastleAppCustomResourcesToForecastleApps",
			args: args{
				forecastleApps: []v1alpha1.ForecastleApp{
					*testutil.CreateForecastleApp("app1", "https://google.com", "default", "https://google.com/icon.png"),
				},
			},
			wantApps: []forecastle.App{
				{
					Name:            "app1",
					Group:           "default",
					Icon:            "https://google.com/icon.png",
					URL:             "https://google.com",
					DiscoverySource: forecastle.ForecastleAppCRD,
				},
			},
		},
		{
			name: "TestConvertForecastleAppCustomResourcesToForecastleAppsWithInvalidAppRouteRef",
			args: args{
				forecastleApps: []v1alpha1.ForecastleApp{
					*testutil.CreateForecastleApp("app1", "https://google.com", "default", "https://google.com/icon.png"),
					*testutil.CreateForecastleAppWithURLFromRoute("invalid-app", "default", "https://google.com/icon.png", "invalid-route"),
				},
			},
			wantApps: []forecastle.App{
				{
					Name:            "app1",
					Group:           "default",
					Icon:            "https://google.com/icon.png",
					URL:             "https://google.com",
					DiscoverySource: forecastle.ForecastleAppCRD,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotApps, err := convertForecastleAppCustomResourcesToForecastleApps(clients, tt.args.forecastleApps); !reflect.DeepEqual(gotApps, tt.wantApps) && err != tt.err {
				t.Errorf("convertForecastleAppCustomResourcesToForecastleApps() = %v, want %v, err = %v, wantErr = %v", gotApps, tt.wantApps, err, tt.err)
			}
		})
	}

}

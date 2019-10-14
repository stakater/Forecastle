package crdapps

import (
	"reflect"
	"testing"

	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	forecastlev1alpha1 "github.com/stakater/Forecastle/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/pkg/client/clientset/versioned/fake"

	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle"
	"github.com/stakater/Forecastle/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewList(t *testing.T) {
	forecastleClient := fake.NewSimpleClientset()

	type args struct {
		forecastleClient forecastlev1alpha1.Interface
		appConfig        config.Config
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithNokubeClient",
			args: args{
				forecastleClient: nil,
			},
			want: &List{
				forecastleClient: nil,
			},
		},
		{
			name: "TestNewListWithDefaultkubeClient",
			args: args{
				forecastleClient: forecastleClient,
			},
			want: &List{
				forecastleClient: forecastleClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.forecastleClient, tt.args.appConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Populate(t *testing.T) {

	forecastleClient := fake.NewSimpleClientset()

	forecastleApp := testutil.CreateForecastleApp("app-1", "https://google.com", "default", "https://google.com/icon.png")

	_, _ = forecastleClient.ForecastleV1alpha1().ForecastleApps("default").Create(forecastleApp)

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
				forecastleClient: forecastleClient,
			},
			args: args{
				namespaces: []string{"default"},
			},
			want: &List{
				forecastleClient: forecastleClient,
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

	_ = forecastleClient.ForecastleV1alpha1().ForecastleApps("default").Delete("app-1", &metav1.DeleteOptions{})

}

func TestList_Get(t *testing.T) {
	forecastleClient := fake.NewSimpleClientset()

	tests := []struct {
		name    string
		al      *List
		want    []forecastle.App
		wantErr bool
	}{
		{
			name: "TestGetForecastleApps",
			al: &List{
				forecastleClient: forecastleClient,
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
				forecastleClient: forecastleClient,
				items:            []forecastle.App{},
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
	type args struct {
		forecastleApps []v1alpha1.ForecastleApp
	}
	tests := []struct {
		name     string
		args     args
		wantApps []forecastle.App
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotApps := convertForecastleAppCustomResourcesToForecastleApps(tt.args.forecastleApps); !reflect.DeepEqual(gotApps, tt.wantApps) {
				t.Errorf("convertForecastleAppCustomResourcesToForecastleApps() = %v, want %v", gotApps, tt.wantApps)
			}
		})
	}

}

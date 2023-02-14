package ingressapps

import (
	"context"
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewList(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type args struct {
		kubeClient kubernetes.Interface
		appConfig  config.Config
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithNokubeClient",
			args: args{
				kubeClient: nil,
			},
			want: &List{
				kubeClient: nil,
			},
		},
		{
			name: "TestNewListWithDefaultkubeClient",
			args: args{
				kubeClient: kubeClient,
			},
			want: &List{
				kubeClient: kubeClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.kubeClient, tt.args.appConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type fields struct {
		kubeClient kubernetes.Interface
		items      []forecastle.App
		err        error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []forecastle.App
		wantErr bool
	}{
		{
			name: "TestGetForecastleApps",
			fields: fields{
				kubeClient: kubeClient,
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
			fields: fields{
				kubeClient: kubeClient,
				items:      []forecastle.App{},
			},
			want:    []forecastle.App{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := &List{
				kubeClient: tt.fields.kubeClient,
				items:      tt.fields.items,
				err:        tt.fields.err,
			}
			got, err := al.Get()
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

func Test_convertIngressesToForecastleApps(t *testing.T) {
	type args struct {
		ingresses []networking.Ingress
	}
	tests := []struct {
		name     string
		args     args
		wantApps []forecastle.App
	}{
		{
			name: "TestConvertIngressesToForecastleAppsWithNoApps",
			args: args{
				ingresses: []networking.Ingress{},
			},
			wantApps: nil,
		},
		{
			name: "TestConvertIngressesToForecastleApps",
			args: args{
				ingresses: []networking.Ingress{
					*testutil.AddAnnotationToIngress(testutil.CreateIngressWithHost("test-ingress", "google.com"),
						annotations.ForecastleIconAnnotation, "https://google.com/icon.png"),
				},
			},
			wantApps: []forecastle.App{
				{
					Name:  "test-ingress",
					Group: "",
					Icon:  "https://google.com/icon.png",
					URL:   "http://google.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotApps := convertIngressesToForecastleApps(tt.args.ingresses); !reflect.DeepEqual(gotApps, tt.wantApps) {
				t.Errorf("convertIngressesToForecastleApps() = %v, want %v", gotApps, tt.wantApps)
			}
		})
	}
}

func TestList_Populate(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	ingress := testutil.AddAnnotationToIngress(
		testutil.AddAnnotationToIngress(
			testutil.CreateIngressWithHost("test-ingress", "google.com"), annotations.ForecastleExposeAnnotation, "true"),
		annotations.ForecastleNetworkRestrictedAnnotation, "true")

	_, _ = kubeClient.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testing"}}, metav1.CreateOptions{})
	_, _ = kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})
	_, _ = kubeClient.NetworkingV1().Ingresses("testing").Create(context.TODO(), ingress, metav1.CreateOptions{})

	type fields struct {
		kubeClient kubernetes.Interface
		items      []forecastle.App
		err        error
	}
	type args struct {
		namespaces []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *List
	}{
		{
			name: "TestListPopulateWithSelectedNamespaces",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{"default", "testing"},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []forecastle.App{
					{
						Name:              "test-ingress",
						Group:             "default",
						URL:               "http://google.com",
						NetworkRestricted: true,
					},
					{
						Name:              "test-ingress",
						Group:             "testing",
						URL:               "http://google.com",
						NetworkRestricted: true,
					},
				},
			},
		},
		{
			name: "TestListPopulateWithAllNamespaces",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{metav1.NamespaceAll},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []forecastle.App{
					{
						Name:              "test-ingress",
						Group:             "default",
						URL:               "http://google.com",
						NetworkRestricted: true,
					},
					{
						Name:              "test-ingress",
						Group:             "testing",
						URL:               "http://google.com",
						NetworkRestricted: true,
					},
				},
			},
		},
		{
			name: "TestListPopulateWithOneNamespace",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{"testing"},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []forecastle.App{
					{
						Name:              "test-ingress",
						Group:             "testing",
						URL:               "http://google.com",
						NetworkRestricted: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := &List{
				kubeClient: tt.fields.kubeClient,
				items:      tt.fields.items,
				err:        tt.fields.err,
			}
			if got := al.Populate(tt.args.namespaces...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Populate() = %v, want %v", got, tt.want)
			}
		})
	}

	_ = kubeClient.CoreV1().Namespaces().Delete(context.TODO(), "testing", metav1.DeleteOptions{})
	_ = kubeClient.NetworkingV1().Ingresses("default").Delete(context.TODO(), "test-ingress", metav1.DeleteOptions{})
	_ = kubeClient.NetworkingV1().Ingresses("testing").Delete(context.TODO(), "test-ingress", metav1.DeleteOptions{})
}

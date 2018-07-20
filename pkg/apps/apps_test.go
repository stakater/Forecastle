package apps

import (
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/pkg/testutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewList(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type args struct {
		kubeClient kubernetes.Interface
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
			if got := NewList(tt.args.kubeClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type fields struct {
		kubeClient kubernetes.Interface
		items      []ForecastleApp
		err        error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ForecastleApp
		wantErr bool
	}{
		{
			name: "TestGetForecastleApps",
			fields: fields{
				kubeClient: kubeClient,
				items: []ForecastleApp{
					{
						Name:      "app",
						Icon:      "https://google.com/icon.png",
						Namespace: "test",
						URL:       "https://google.com",
					},
				},
			},
			want: []ForecastleApp{
				{
					Name:      "app",
					Icon:      "https://google.com/icon.png",
					Namespace: "test",
					URL:       "https://google.com",
				},
			},
			wantErr: false,
		},
		{
			name: "TestGetForecastleAppsWithEmptyList",
			fields: fields{
				kubeClient: kubeClient,
				items:      []ForecastleApp{},
			},
			want:    []ForecastleApp{},
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
		ingresses []v1beta1.Ingress
	}
	tests := []struct {
		name     string
		args     args
		wantApps []ForecastleApp
	}{
		{
			name: "TestConvertIngressesToForecastleAppsWithNoApps",
			args: args{
				ingresses: []v1beta1.Ingress{},
			},
			wantApps: nil,
		},
		{
			name: "TestConvertIngressesToForecastleApps",
			args: args{
				ingresses: []v1beta1.Ingress{
					*testutil.AddAnnotationToIngress(testutil.CreateIngressWithHost("test-ingress", "google.com"),
						ForecastleIconAnnotation, "https://google.com/icon.png"),
				},
			},
			wantApps: []ForecastleApp{
				{
					Name:      "test-ingress",
					Namespace: "",
					Icon:      "https://google.com/icon.png",
					URL:       "http://google.com",
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
			testutil.CreateIngressWithHost("test-ingress", "google.com"), ForecastleExposeAnnotation, "true"),
		IngressClassAnnotation, "ingress")

	_, _ = kubeClient.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testing"}})
	_, _ = kubeClient.ExtensionsV1beta1().Ingresses("default").Create(ingress)
	_, _ = kubeClient.ExtensionsV1beta1().Ingresses("testing").Create(ingress)

	type fields struct {
		kubeClient kubernetes.Interface
		items      []ForecastleApp
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
				items: []ForecastleApp{
					{
						Name:      "test-ingress",
						Namespace: "default",
						URL:       "http://google.com",
					},
					{
						Name:      "test-ingress",
						Namespace: "testing",
						URL:       "http://google.com",
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
				items: []ForecastleApp{
					{
						Name:      "test-ingress",
						Namespace: "default",
						URL:       "http://google.com",
					},
					{
						Name:      "test-ingress",
						Namespace: "testing",
						URL:       "http://google.com",
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
				items: []ForecastleApp{
					{
						Name:      "test-ingress",
						Namespace: "testing",
						URL:       "http://google.com",
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

	_ = kubeClient.CoreV1().Namespaces().Delete("testing", &metav1.DeleteOptions{})
	_ = kubeClient.ExtensionsV1beta1().Ingresses("default").Delete("test-ingress", &metav1.DeleteOptions{})
	_ = kubeClient.ExtensionsV1beta1().Ingresses("testing").Delete("test-ingress", &metav1.DeleteOptions{})
}

package ingresses

import (
	"context"
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/config"
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

func TestList_Populate(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	ingress := testutil.CreateIngressWithHost("test-ingress", "google.com")

	_, _ = kubeClient.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testing"}}, metav1.CreateOptions{})
	ingressDefault, _ := kubeClient.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})
	ingressTesting, _ := kubeClient.NetworkingV1().Ingresses("testing").Create(context.TODO(), ingress, metav1.CreateOptions{})

	type fields struct {
		items      []networking.Ingress
		err        error
		kubeClient kubernetes.Interface
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
			name: "TestPopulateWithSelectedNamespaces",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{"default", "testing"},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []networking.Ingress{
					*ingressDefault,
					*ingressTesting,
				},
			},
		},
		{
			name: "TestPopulateWithAllNamespaces",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{metav1.NamespaceAll},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []networking.Ingress{
					*ingressDefault,
					*ingressTesting,
				},
			},
		},
		{
			name: "TestPopulateWithOneNamespace",
			fields: fields{
				kubeClient: kubeClient,
			},
			args: args{
				namespaces: []string{"default"},
			},
			want: &List{
				kubeClient: kubeClient,
				items: []networking.Ingress{
					*ingressDefault,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &List{
				items:      tt.fields.items,
				err:        tt.fields.err,
				kubeClient: tt.fields.kubeClient,
			}
			if got := il.Populate(tt.args.namespaces...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Populate() = %v, want %v", got, tt.want)
			}
		})
	}

	_ = kubeClient.CoreV1().Namespaces().Delete(context.TODO(), "testing", metav1.DeleteOptions{})
	_ = kubeClient.NetworkingV1().Ingresses("default").Delete(context.TODO(), "test-ingress", metav1.DeleteOptions{})
	_ = kubeClient.NetworkingV1().Ingresses("testing").Delete(context.TODO(), "test-ingress", metav1.DeleteOptions{})
}

func TestList_Filter(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	type fields struct {
		items      []networking.Ingress
		err        error
		kubeClient kubernetes.Interface
	}
	type args struct {
		filterFunc func(networking.Ingress, config.Config) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *List
	}{
		{
			name: "TestListFilter",
			fields: fields{
				items: []networking.Ingress{
					*testutil.CreateIngress("test-ingress"),
				},
				kubeClient: kubeClient,
			},
			args: args{
				filterFunc: func(ingress networking.Ingress, appConfig config.Config) bool { return true },
			},
			want: &List{
				kubeClient: kubeClient,
				items: []networking.Ingress{
					*testutil.CreateIngress("test-ingress"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &List{
				items:      tt.fields.items,
				err:        tt.fields.err,
				kubeClient: tt.fields.kubeClient,
			}
			if got := il.Filter(tt.args.filterFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()

	type fields struct {
		items      []networking.Ingress
		err        error
		kubeClient kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []networking.Ingress
		wantErr bool
	}{
		{
			name: "TestListGetWithNoItems",
			fields: fields{
				kubeClient: kubeClient,
			},
			want:    nil,
			wantErr: false,
		}, {
			name: "TestListGetWithItems",
			fields: fields{
				kubeClient: kubeClient,
				items: []networking.Ingress{
					*testutil.CreateIngress("test-ingress"),
				},
			},
			want: []networking.Ingress{
				*testutil.CreateIngress("test-ingress"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &List{
				items:      tt.fields.items,
				err:        tt.fields.err,
				kubeClient: tt.fields.kubeClient,
			}
			got, err := il.Get()
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

package httproutes

import (
	"context"
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gateway "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
	gatewayfake "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/fake"
)

func TestNewList(t *testing.T) {
	gatewayClient := gatewayfake.NewSimpleClientset()
	type args struct {
		gatewayClient gateway.Interface
		appConfig     config.Config
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithNoClient",
			args: args{
				gatewayClient: nil,
			},
			want: &List{
				gatewayClient: nil,
			},
		},
		{
			name: "TestNewListWithClient",
			args: args{
				gatewayClient: gatewayClient,
			},
			want: &List{
				gatewayClient: gatewayClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewList(tt.args.gatewayClient, tt.args.appConfig); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewList() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestList_Populate(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() //nolint:staticcheck // NewClientset requires generated apply configurations
	gatewayClient := gatewayfake.NewSimpleClientset()

	httpRoute := testutil.CreateHTTPRouteWithHostnameAndNamespace("test-route", "default", "app.example.com")

	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testing"}}, metav1.CreateOptions{},
	)
	httpRouteDefault, _ := gatewayClient.GatewayV1().HTTPRoutes("default").Create(context.TODO(), httpRoute, metav1.CreateOptions{})

	httpRouteTesting := testutil.CreateHTTPRouteWithHostnameAndNamespace("test-route", "testing", "app.example.com")
	httpRouteTesting, _ = gatewayClient.GatewayV1().HTTPRoutes("testing").Create(context.TODO(), httpRouteTesting, metav1.CreateOptions{})

	type fields struct {
		items         []gatewayv1.HTTPRoute
		err           error
		gatewayClient gateway.Interface
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
				gatewayClient: gatewayClient,
			},
			args: args{
				namespaces: []string{"default", "testing"},
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []gatewayv1.HTTPRoute{
					*httpRouteDefault,
					*httpRouteTesting,
				},
			},
		},
		{
			name: "TestPopulateWithAllNamespaces",
			fields: fields{
				gatewayClient: gatewayClient,
			},
			args: args{
				namespaces: []string{metav1.NamespaceAll},
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []gatewayv1.HTTPRoute{
					*httpRouteDefault,
					*httpRouteTesting,
				},
			},
		},
		{
			name: "TestPopulateWithOneNamespace",
			fields: fields{
				gatewayClient: gatewayClient,
			},
			args: args{
				namespaces: []string{"default"},
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []gatewayv1.HTTPRoute{
					*httpRouteDefault,
				},
			},
		},
		{
			name: "TestPopulateWithNilClient",
			fields: fields{
				gatewayClient: nil,
			},
			args: args{
				namespaces: []string{"default"},
			},
			want: &List{
				gatewayClient: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hl := &List{
					items:         tt.fields.items,
					err:           tt.fields.err,
					gatewayClient: tt.fields.gatewayClient,
				}
				if got := hl.Populate(tt.args.namespaces...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List.Populate() = %v, want %v", got, tt.want)
				}
			},
		)
	}

	_ = kubeClient.CoreV1().Namespaces().Delete(context.TODO(), "testing", metav1.DeleteOptions{})
	_ = gatewayClient.GatewayV1().HTTPRoutes("default").Delete(context.TODO(), "test-route", metav1.DeleteOptions{})
	_ = gatewayClient.GatewayV1().HTTPRoutes("testing").Delete(context.TODO(), "test-route", metav1.DeleteOptions{})
}

func TestList_Filter(t *testing.T) {
	gatewayClient := gatewayfake.NewSimpleClientset()

	type fields struct {
		items         []gatewayv1.HTTPRoute
		err           error
		gatewayClient gateway.Interface
	}
	type args struct {
		filterFunc func(gatewayv1.HTTPRoute, config.Config) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *List
	}{
		{
			name: "TestListFilterReturnsAll",
			fields: fields{
				items: []gatewayv1.HTTPRoute{
					*testutil.CreateHTTPRoute("test-route"),
				},
				gatewayClient: gatewayClient,
			},
			args: args{
				filterFunc: func(httpRoute gatewayv1.HTTPRoute, appConfig config.Config) bool { return true },
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []gatewayv1.HTTPRoute{
					*testutil.CreateHTTPRoute("test-route"),
				},
			},
		},
		{
			name: "TestListFilterReturnsNone",
			fields: fields{
				items: []gatewayv1.HTTPRoute{
					*testutil.CreateHTTPRoute("test-route"),
				},
				gatewayClient: gatewayClient,
			},
			args: args{
				filterFunc: func(httpRoute gatewayv1.HTTPRoute, appConfig config.Config) bool { return false },
			},
			want: &List{
				gatewayClient: gatewayClient,
				items:         nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hl := &List{
					items:         tt.fields.items,
					err:           tt.fields.err,
					gatewayClient: tt.fields.gatewayClient,
				}
				if got := hl.Filter(tt.args.filterFunc); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List.Filter() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestList_Get(t *testing.T) {
	gatewayClient := gatewayfake.NewSimpleClientset()

	type fields struct {
		items         []gatewayv1.HTTPRoute
		err           error
		gatewayClient gateway.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []gatewayv1.HTTPRoute
		wantErr bool
	}{
		{
			name: "TestListGetWithNoItems",
			fields: fields{
				gatewayClient: gatewayClient,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "TestListGetWithItems",
			fields: fields{
				gatewayClient: gatewayClient,
				items: []gatewayv1.HTTPRoute{
					*testutil.CreateHTTPRoute("test-route"),
				},
			},
			want: []gatewayv1.HTTPRoute{
				*testutil.CreateHTTPRoute("test-route"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hl := &List{
					items:         tt.fields.items,
					err:           tt.fields.err,
					gatewayClient: tt.fields.gatewayClient,
				}
				got, err := hl.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("List.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List.Get() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

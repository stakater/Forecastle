package httprouteapps

import (
	"context"
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
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

func TestList_Get(t *testing.T) {
	gatewayClient := gatewayfake.NewSimpleClientset()
	type fields struct {
		gatewayClient gateway.Interface
		items         []forecastle.App
		err           error
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
				gatewayClient: gatewayClient,
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
				gatewayClient: gatewayClient,
				items:         []forecastle.App{},
			},
			want:    []forecastle.App{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				al := &List{
					gatewayClient: tt.fields.gatewayClient,
					items:         tt.fields.items,
					err:           tt.fields.err,
				}
				got, err := al.Get()
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

func Test_convertHTTPRoutesToForecastleApps(t *testing.T) {
	type args struct {
		httpRoutes []gatewayv1.HTTPRoute
	}
	tests := []struct {
		name     string
		args     args
		wantApps []forecastle.App
	}{
		{
			name: "TestConvertHTTPRoutesToForecastleAppsWithNoApps",
			args: args{
				httpRoutes: []gatewayv1.HTTPRoute{},
			},
			wantApps: nil,
		},
		{
			name: "TestConvertHTTPRoutesToForecastleApps",
			args: args{
				httpRoutes: []gatewayv1.HTTPRoute{
					*testutil.AddAnnotationToHTTPRoute(
						testutil.CreateHTTPRouteWithHostname("test-route", "app.example.com"),
						annotations.ForecastleIconAnnotation, "https://example.com/icon.png",
					),
				},
			},
			wantApps: []forecastle.App{
				{
					Name:            "test-route",
					Group:           "",
					Icon:            "https://example.com/icon.png",
					URL:             "https://app.example.com",
					DiscoverySource: forecastle.HTTPRoute,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotApps := convertHTTPRoutesToForecastleApps(tt.args.httpRoutes); !reflect.DeepEqual(gotApps, tt.wantApps) {
					t.Errorf("convertHTTPRoutesToForecastleApps() = %v, want %v", gotApps, tt.wantApps)
				}
			},
		)
	}
}

func TestList_Populate(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() //nolint:staticcheck // NewClientset requires generated apply configurations
	gatewayClient := gatewayfake.NewSimpleClientset()

	httpRoute := testutil.AddAnnotationToHTTPRoute(
		testutil.AddAnnotationToHTTPRoute(
			testutil.CreateHTTPRouteWithHostnameAndNamespace("test-route", "default", "app.example.com"),
			annotations.ForecastleExposeAnnotation, "true",
		),
		annotations.ForecastleNetworkRestrictedAnnotation, "true",
	)

	_, _ = kubeClient.CoreV1().Namespaces().Create(
		context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testing"}}, metav1.CreateOptions{},
	)
	_, _ = gatewayClient.GatewayV1().HTTPRoutes("default").Create(context.TODO(), httpRoute, metav1.CreateOptions{})

	httpRouteTesting := testutil.AddAnnotationToHTTPRoute(
		testutil.AddAnnotationToHTTPRoute(
			testutil.CreateHTTPRouteWithHostnameAndNamespace("test-route", "testing", "app.example.com"),
			annotations.ForecastleExposeAnnotation, "true",
		),
		annotations.ForecastleNetworkRestrictedAnnotation, "true",
	)
	_, _ = gatewayClient.GatewayV1().HTTPRoutes("testing").Create(context.TODO(), httpRouteTesting, metav1.CreateOptions{})

	type fields struct {
		gatewayClient gateway.Interface
		items         []forecastle.App
		err           error
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
				gatewayClient: gatewayClient,
			},
			args: args{
				namespaces: []string{"default", "testing"},
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []forecastle.App{
					{
						Name:              "test-route",
						Group:             "default",
						URL:               "https://app.example.com",
						DiscoverySource:   forecastle.HTTPRoute,
						NetworkRestricted: true,
					},
					{
						Name:              "test-route",
						Group:             "testing",
						URL:               "https://app.example.com",
						DiscoverySource:   forecastle.HTTPRoute,
						NetworkRestricted: true,
					},
				},
			},
		},
		{
			name: "TestListPopulateWithOneNamespace",
			fields: fields{
				gatewayClient: gatewayClient,
			},
			args: args{
				namespaces: []string{"testing"},
			},
			want: &List{
				gatewayClient: gatewayClient,
				items: []forecastle.App{
					{
						Name:              "test-route",
						Group:             "testing",
						URL:               "https://app.example.com",
						DiscoverySource:   forecastle.HTTPRoute,
						NetworkRestricted: true,
					},
				},
			},
		},
		{
			name: "TestListPopulateWithNilClient",
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
				al := &List{
					gatewayClient: tt.fields.gatewayClient,
					items:         tt.fields.items,
					err:           tt.fields.err,
				}
				if got := al.Populate(tt.args.namespaces...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List.Populate() = %v, want %v", got, tt.want)
				}
			},
		)
	}

	_ = kubeClient.CoreV1().Namespaces().Delete(context.TODO(), "testing", metav1.DeleteOptions{})
	_ = gatewayClient.GatewayV1().HTTPRoutes("default").Delete(context.TODO(), "test-route", metav1.DeleteOptions{})
	_ = gatewayClient.GatewayV1().HTTPRoutes("testing").Delete(context.TODO(), "test-route", metav1.DeleteOptions{})
}

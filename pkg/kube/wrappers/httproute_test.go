package wrappers

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestHTTPRouteWrapper_GetURL(t *testing.T) {
	tests := []struct {
		name      string
		httpRoute *gatewayv1.HTTPRoute
		want      string
	}{
		{
			name: "WithSingleHostname",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "test-route"},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{"app.example.com"},
				},
			},
			want: "https://app.example.com",
		},
		{
			name: "WithMultipleHostnames",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "test-route"},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{"primary.example.com", "secondary.example.com"},
				},
			},
			want: "https://primary.example.com",
		},
		{
			name: "WithNoHostnames",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "test-route"},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{},
				},
			},
			want: "",
		},
		{
			name: "WithNilHostnames",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "test-route"},
				Spec:       gatewayv1.HTTPRouteSpec{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			if got := hw.GetURL(); got != tt.want {
				t.Errorf("HTTPRouteWrapper.GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

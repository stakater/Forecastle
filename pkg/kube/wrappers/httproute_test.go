package wrappers

import (
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
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
			name:      "WithSingleHostname",
			httpRoute: testutil.CreateHTTPRouteWithHostname("test-route", "app.example.com"),
			want:      "https://app.example.com",
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
			name:      "WithNoHostnames",
			httpRoute: testutil.CreateHTTPRoute("test-route"),
			want:      "",
		},
		{
			name: "WithURLAnnotation",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRouteWithHostname("test-route", "app.example.com"),
				annotations.ForecastleURLAnnotation, "https://logging.example.net/select/vmui/"),
			want: "https://logging.example.net/select/vmui/",
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

func TestHTTPRouteWrapper_GetName(t *testing.T) {
	tests := []struct {
		name      string
		httpRoute *gatewayv1.HTTPRoute
		want      string
	}{
		{
			name:      "WithResourceName",
			httpRoute: testutil.CreateHTTPRoute("my-route"),
			want:      "my-route",
		},
		{
			name: "WithAppNameAnnotation",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRoute("my-route"),
				annotations.ForecastleAppNameAnnotation, "My App"),
			want: "My App",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			if got := hw.GetName(); got != tt.want {
				t.Errorf("HTTPRouteWrapper.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRouteWrapper_GetNamespace(t *testing.T) {
	tests := []struct {
		name      string
		httpRoute *gatewayv1.HTTPRoute
		want      string
	}{
		{
			name:      "WithNamespace",
			httpRoute: testutil.CreateHTTPRouteWithNamespace("my-route", "production"),
			want:      "production",
		},
		{
			name:      "WithEmptyNamespace",
			httpRoute: testutil.CreateHTTPRoute("my-route"),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			if got := hw.GetNamespace(); got != tt.want {
				t.Errorf("HTTPRouteWrapper.GetNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRouteWrapper_GetGroup(t *testing.T) {
	tests := []struct {
		name      string
		httpRoute *gatewayv1.HTTPRoute
		want      string
	}{
		{
			name:      "DefaultsToNamespace",
			httpRoute: testutil.CreateHTTPRouteWithNamespace("my-route", "Production"),
			want:      "production",
		},
		{
			name: "WithGroupAnnotation",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRouteWithNamespace("my-route", "default"),
				annotations.ForecastleGroupAnnotation, "My Group"),
			want: "my group",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			if got := hw.GetGroup(); got != tt.want {
				t.Errorf("HTTPRouteWrapper.GetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRouteWrapper_GetAnnotationValue(t *testing.T) {
	tests := []struct {
		name          string
		httpRoute     *gatewayv1.HTTPRoute
		annotationKey string
		want          string
	}{
		{
			name:          "WithNoAnnotations",
			httpRoute:     testutil.CreateHTTPRoute("my-route"),
			annotationKey: "some-key",
			want:          "",
		},
		{
			name: "WithMatchingAnnotation",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRoute("my-route"),
				annotations.ForecastleIconAnnotation, "https://example.com/icon.png"),
			annotationKey: annotations.ForecastleIconAnnotation,
			want:          "https://example.com/icon.png",
		},
		{
			name: "WithNonMatchingAnnotation",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRoute("my-route"),
				annotations.ForecastleIconAnnotation, "https://example.com/icon.png"),
			annotationKey: annotations.ForecastleGroupAnnotation,
			want:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			if got := hw.GetAnnotationValue(tt.annotationKey); got != tt.want {
				t.Errorf("HTTPRouteWrapper.GetAnnotationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRouteWrapper_GetProperties(t *testing.T) {
	tests := []struct {
		name      string
		httpRoute *gatewayv1.HTTPRoute
		want      map[string]string
	}{
		{
			name:      "WithNoProperties",
			httpRoute: testutil.CreateHTTPRoute("my-route"),
			want:      nil,
		},
		{
			name: "WithProperties",
			httpRoute: testutil.AddAnnotationToHTTPRoute(
				testutil.CreateHTTPRoute("my-route"),
				annotations.ForecastlePropertiesAnnotation, "key1:value1,key2:value2"),
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := NewHTTPRouteWrapper(tt.httpRoute)
			got := hw.GetProperties()
			if len(got) != len(tt.want) {
				t.Errorf("HTTPRouteWrapper.GetProperties() = %v, want %v", got, tt.want)
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("HTTPRouteWrapper.GetProperties()[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

package apps

import (
	"testing"

	"github.com/stakater/Forecastle/pkg/testutil"
	"k8s.io/api/extensions/v1beta1"
)

func Test_byIngressClassAnnotation(t *testing.T) {
	type args struct {
		ingress v1beta1.Ingress
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestByIngressClassAnnotationWithNoAnnotation",
			args: args{
				ingress: *testutil.CreateIngress("test-ingress"),
			},
			want: false,
		},
		{
			name: "TestByIngressClassAnnotationWithAnnotation",
			args: args{
				ingress: *testutil.AddAnnotationToIngress(testutil.CreateIngress("test-ingress"), IngressClassAnnotation, "internal-ingress"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byIngressClassAnnotation(tt.args.ingress); got != tt.want {
				t.Errorf("byIngressClassAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byForecastleExposeAnnotation(t *testing.T) {
	type args struct {
		ingress v1beta1.Ingress
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestByForecastleExposeAnnotationWithNoAnnotation",
			args: args{
				ingress: *testutil.CreateIngress("test-ingress"),
			},
			want: false,
		},
		{
			name: "TestByForecastleExposeAnnotationWithAnnotationsFalseValue",
			args: args{
				ingress: *testutil.AddAnnotationToIngress(testutil.CreateIngress("test-ingress"), ForecastleExposeAnnotation, "false"),
			},
			want: false,
		},
		{
			name: "TestByForecastleExposeAnnotationWithAnnotationsTrueValue",
			args: args{
				ingress: *testutil.AddAnnotationToIngress(testutil.CreateIngress("test-ingress"), ForecastleExposeAnnotation, "true"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byForecastleExposeAnnotation(tt.args.ingress); got != tt.want {
				t.Errorf("byForecastleExposeAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

package ingressapps

import (
	"testing"

	"github.com/stakater/Forecastle/pkg/annotations"
	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/testutil"
	v1 "k8s.io/api/networking/v1"
)

func Test_byForecastleExposeAnnotation(t *testing.T) {
	type args struct {
		ingress   v1.Ingress
		appConfig config.Config
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
				ingress: *testutil.AddAnnotationToIngress(
					testutil.CreateIngress("test-ingress"), annotations.ForecastleExposeAnnotation, "false"),
			},
			want: false,
		},
		{
			name: "TestByForecastleExposeAnnotationWithTrueValue",
			args: args{
				ingress: *testutil.AddAnnotationToIngress(
					testutil.CreateIngress("test-ingress"), annotations.ForecastleExposeAnnotation, "true"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byForecastleExposeAnnotation(tt.args.ingress, tt.args.appConfig); got != tt.want {
				t.Errorf("byForecastleExposeAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

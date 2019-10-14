package crdapps

import (
	"testing"

	v1alpha1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/pkg/config"

	"github.com/stakater/Forecastle/pkg/testutil"
)

func Test_byForecastleInstance(t *testing.T) {
	app := testutil.CreateForecastleApp("app-1", "https://google.com", "default", "https://google.com/icon.png")
	app.Spec.Instance = "default"

	type args struct {
		forecastleApp v1alpha1.ForecastleApp
		appConfig     config.Config
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestByForecastleInstanceWithMatchingInstance",
			args: args{
				forecastleApp: *app,
				appConfig: config.Config{
					InstanceName: "default",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byForecastleInstance(tt.args.forecastleApp, tt.args.appConfig); got != tt.want {
				t.Errorf("byForecastleInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

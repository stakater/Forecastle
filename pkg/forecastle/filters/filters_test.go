package filters

import (
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/config"
)

func TestByForecastleExposeAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		annots    map[string]string
		appConfig config.Config
		want      bool
	}{
		{
			name:   "NoAnnotations",
			annots: nil,
			want:   false,
		},
		{
			name:   "EmptyAnnotations",
			annots: map[string]string{},
			want:   false,
		},
		{
			name: "ExposeAnnotationFalse",
			annots: map[string]string{
				annotations.ForecastleExposeAnnotation: "false",
			},
			want: false,
		},
		{
			name: "ExposeAnnotationTrue",
			annots: map[string]string{
				annotations.ForecastleExposeAnnotation: "true",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByForecastleExposeAnnotation(tt.annots, tt.appConfig); got != tt.want {
				t.Errorf("ByForecastleExposeAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByForecastleInstanceAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		annots    map[string]string
		appConfig config.Config
		want      bool
	}{
		{
			name:      "NoAnnotations",
			annots:    nil,
			appConfig: config.Config{InstanceName: "test"},
			want:      false,
		},
		{
			name: "InstanceMatches",
			annots: map[string]string{
				annotations.ForecastleInstanceAnnotation: "test",
			},
			appConfig: config.Config{InstanceName: "test"},
			want:      true,
		},
		{
			name: "InstanceInList",
			annots: map[string]string{
				annotations.ForecastleInstanceAnnotation: "foo,test,bar",
			},
			appConfig: config.Config{InstanceName: "test"},
			want:      true,
		},
		{
			name: "InstanceNotInList",
			annots: map[string]string{
				annotations.ForecastleInstanceAnnotation: "foo,bar",
			},
			appConfig: config.Config{InstanceName: "test"},
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByForecastleInstanceAnnotation(tt.annots, tt.appConfig); got != tt.want {
				t.Errorf("ByForecastleInstanceAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByInstance(t *testing.T) {
	tests := []struct {
		name          string
		instanceValue string
		appConfig     config.Config
		want          bool
	}{
		{
			name:          "EmptyInstanceValue",
			instanceValue: "",
			appConfig:     config.Config{InstanceName: "test"},
			want:          false,
		},
		{
			name:          "InstanceMatches",
			instanceValue: "test",
			appConfig:     config.Config{InstanceName: "test"},
			want:          true,
		},
		{
			name:          "InstanceInList",
			instanceValue: "foo,test,bar",
			appConfig:     config.Config{InstanceName: "test"},
			want:          true,
		},
		{
			name:          "InstanceNotInList",
			instanceValue: "foo,bar",
			appConfig:     config.Config{InstanceName: "test"},
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByInstance(tt.instanceValue, tt.appConfig); got != tt.want {
				t.Errorf("ByInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

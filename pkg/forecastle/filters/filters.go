package filters

import (
	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/util/strings"
)

// ByForecastleExposeAnnotation returns true if annotations have forecastle expose set to true
func ByForecastleExposeAnnotation(annots map[string]string, appConfig config.Config) bool {
	if annots != nil {
		if val, ok := annots[annotations.ForecastleExposeAnnotation]; ok {
			return val == "true"
		}
	}
	return false
}

// ByForecastleInstanceAnnotation returns true if annotations match the configured forecastle instance
func ByForecastleInstanceAnnotation(annots map[string]string, appConfig config.Config) bool {
	if annots != nil {
		if val, ok := annots[annotations.ForecastleInstanceAnnotation]; ok {
			return ByInstance(val, appConfig)
		}
	}
	return false
}

// ByInstance returns true if the instance value matches the configured forecastle instance
func ByInstance(instanceValue string, appConfig config.Config) bool {
	return strings.ContainsBetweenDelimiter(instanceValue, appConfig.InstanceName, ",")
}

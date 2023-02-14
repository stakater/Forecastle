package crdapps

import (
	v1alpha1 "github.com/stakater/Forecastle/v1/pkg/apis/forecastle/v1alpha1"
	"github.com/stakater/Forecastle/v1/pkg/config"

	"github.com/stakater/Forecastle/v1/pkg/util/strings"
)

func byForecastleInstance(forecastleApp v1alpha1.ForecastleApp, appConfig config.Config) bool {
	return strings.ContainsBetweenDelimiter(forecastleApp.Spec.Instance, appConfig.InstanceName, ",")
}

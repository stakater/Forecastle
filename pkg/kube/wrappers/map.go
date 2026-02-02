package wrappers

import (
	"net/url"
	"strings"
)

func makeMap(value string) map[string]string {
	propertiesMap := make(map[string]string)

	propertyParams := strings.Split(value, ",")
	for _, propertyParam := range propertyParams {
		splitted := strings.SplitN(propertyParam, ":", 2)
		propertiesMap[splitted[0]] = splitted[1]
	}

	return propertiesMap
}

func getAnnotationValue(annotations map[string]string, key string) string {
	if annotations != nil {
		if value, ok := annotations[key]; ok {
			return value
		}
	}
	return ""
}

func getAndValidateURLAnnotation(annotations map[string]string, key string) string {
	urlValue := getAnnotationValue(annotations, key)
	if urlValue == "" {
		return ""
	}

	parsedURL, err := url.ParseRequestURI(urlValue)
	if err != nil {
		logger.Warn(err)
		return ""
	}

	return parsedURL.String()
}

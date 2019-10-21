package wrappers

import "strings"

func makeMap(value string) map[string]string {
	detailsMap := make(map[string]string)

	detailParams := strings.Split(value, ",")
	for _, detailParam := range detailParams {
		splitted := strings.SplitN(detailParam, ":", 2)
		detailsMap[splitted[0]] = splitted[1]
	}

	return detailsMap
}

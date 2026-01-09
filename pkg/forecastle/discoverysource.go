package forecastle

import (
	"encoding/json"
	"fmt"
)

type DiscoverySource int

const (
	Ingress DiscoverySource = iota
	Config
	ForecastleAppCRD
	HTTPRoute
)

func (ds DiscoverySource) String() string {
	names := [...]string{
		"Ingress",
		"Config",
		"ForecastleAppCRD",
		"HTTPRoute",
	}

	if ds < Ingress || ds > HTTPRoute {
		return "Unknown"
	}

	return names[ds]
}

func (ds DiscoverySource) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ds.String() + `"`), nil
}

func (ds *DiscoverySource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "Ingress":
		*ds = Ingress
	case "Config":
		*ds = Config
	case "ForecastleAppCRD":
		*ds = ForecastleAppCRD
	case "HTTPRoute":
		*ds = HTTPRoute
	default:
		return fmt.Errorf("unknown DiscoverySource: %s", s)
	}

	return nil
}

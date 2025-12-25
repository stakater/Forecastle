package forecastle

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

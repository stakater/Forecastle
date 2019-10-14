package forecastle

type DiscoverySource int

const (
	Ingress DiscoverySource = iota
	Config
	ForecastleAppCRD
)

func (ds DiscoverySource) String() string {
	names := [...]string{
		"Ingress",
		"Config",
		"ForecastleAppCRD",
	}

	if ds < Ingress || ds > ForecastleAppCRD {
		return "Unknown"
	}

	return names[ds]
}

func (ds DiscoverySource) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ds.String() + `"`), nil
}

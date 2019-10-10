package forecastle

type DiscoverySource int

const (
	Ingress DiscoverySource = iota
	Config
	ForecastleApp
)

func (ds DiscoverySource) String() string {
	names := [...]string{
		"Ingress",
		"Config",
		"ForecastleApp",
	}

	if ds < Ingress || ds > ForecastleApp {
		return "Unknown"
	}

	return names[ds]
}

func (ds DiscoverySource) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ds.String() + `"`), nil
}

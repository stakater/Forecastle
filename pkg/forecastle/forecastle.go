package forecastle

// App struct that contains information about an app that is exposed to forecastle
type App struct {
	Name              string            `json:"name"`
	Icon              string            `json:"icon"`
	Group             string            `json:"group"`
	URL               string            `json:"url"`
	DiscoverySource   DiscoverySource   `json:"discoverySource"`
	NetworkRestricted bool              `json:"networkRestricted"`
	Properties        map[string]string `json:"properties"`
}

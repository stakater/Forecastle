package forecastle

// App struct that contains information about an app that is exposed to forecastle
type App struct {
	Name              string
	Icon              string
	Group             string
	URL               string
	DiscoverySource   DiscoverySource
	NetworkRestricted bool
}

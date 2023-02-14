package customapps

import (
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
)

// un-used variable
// var (
// 	logger = log.New()
// )

// List struct is used for listing forecastle apps
type List struct {
	appConfig config.Config
	err       error // Used for forwarding errors
	items     []forecastle.App
}

// NewList func creates a new instance of apps lister
func NewList(appConfig config.Config) *List {
	return &List{
		appConfig: appConfig,
	}
}

// Populate function that populates a list of custom apps
func (al *List) Populate() *List {
	al.items = convertCustomAppsToForecastleApps(al.appConfig.CustomApps)

	return al
}

// Get function returns the apps currently present in List
func (al *List) Get() ([]forecastle.App, error) {
	return al.items, al.err
}

func convertCustomAppsToForecastleApps(customApps []config.CustomApp) (apps []forecastle.App) {
	for _, customApp := range customApps {
		apps = append(apps, forecastle.App{
			Name:              customApp.Name,
			URL:               customApp.URL,
			Icon:              customApp.Icon,
			Group:             customApp.Group,
			DiscoverySource:   forecastle.Config,
			NetworkRestricted: customApp.NetworkRestricted,
			Properties:        customApp.Properties,
		})
	}

	return apps
}

package forecastlecrdapps

import (
	v1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1"
	forecastlev1 "github.com/stakater/Forecastle/pkg/client/clientset/versioned"
	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle"
	"github.com/stakater/Forecastle/pkg/kube/lists/forecastleapps"
	"github.com/stakater/Forecastle/pkg/log"
)

var (
	logger = log.New()
)

// List struct is used for listing forecastle apps
type List struct {
	appConfig        config.Config
	err              error // Used for forwarding errors
	items            []forecastle.App
	forecastleClient forecastlev1.Interface
}

// NewList func creates a new instance of apps lister
func NewList(forecastleClient forecastlev1.Interface, appConfig config.Config) *List {
	return &List{
		appConfig:        appConfig,
		forecastleClient: forecastleClient,
	}
}

// Populate function that populates a list of forecastle apps from forecastleapps in selected namespaces
func (al *List) Populate(namespaces ...string) *List {
	forecastleAppListObj := forecastleapps.NewList(al.forecastleClient, al.appConfig).
		Populate(namespaces...)

	var forecastleAppList []v1.ForecastleApp
	var err error

	// Apply Instance filter
	if len(al.appConfig.InstanceName) != 0 {
		forecastleAppList, err = forecastleAppListObj.
			Filter(byForecastleInstance).
			Get()
	} else {
		forecastleAppList, err = forecastleAppListObj.Get()
	}

	if err != nil {
		al.err = err
	}

	al.items = convertForecastleAppCustomResourcesToForecastleApps(forecastleAppList)

	return al
}

// Get function returns the apps currently present in List
func (al *List) Get() ([]forecastle.App, error) {
	return al.items, al.err
}

func convertForecastleAppCustomResourcesToForecastleApps(forecastleApps []v1.ForecastleApp) (apps []forecastle.App) {
	for _, forecastleApp := range forecastleApps {
		logger.Infof("Found forecastleApp with Name '%v' in Namespace '%v'", forecastleApp.Name, forecastleApp.Namespace)

		apps = append(apps, forecastle.App{
			Name:            forecastleApp.Spec.Name,
			Group:           forecastleApp.Spec.Group,
			Icon:            forecastleApp.Spec.Icon,
			URL:             forecastleApp.Spec.URL,
			DiscoverySource: forecastle.ForecastleApp,
		})
	}
	return
}

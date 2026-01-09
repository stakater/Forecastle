package httprouteapps

import (
	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/filters"
	"github.com/stakater/Forecastle/v1/pkg/kube/lists/httproutes"
	"github.com/stakater/Forecastle/v1/pkg/kube/wrappers"
	"github.com/stakater/Forecastle/v1/pkg/log"
	"github.com/stakater/Forecastle/v1/pkg/util/strings"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gateway "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

var logger = log.New()

// List struct is used for listing forecastle apps from HTTPRoutes
type List struct {
	appConfig     config.Config
	err           error
	items         []forecastle.App
	gatewayClient gateway.Interface
}

// NewList creates a new instance of apps lister for HTTPRoutes
func NewList(gatewayClient gateway.Interface, appConfig config.Config) *List {
	return &List{
		appConfig:     appConfig,
		gatewayClient: gatewayClient,
	}
}

// Populate populates a list of forecastle apps from HTTPRoutes in selected namespaces
func (al *List) Populate(namespaces ...string) *List {
	if al.gatewayClient == nil {
		return al
	}

	httpRouteList, err := httproutes.NewList(al.gatewayClient, al.appConfig).
		Populate(namespaces...).
		Filter(func(hr gatewayv1.HTTPRoute, cfg config.Config) bool {
			return filters.ByForecastleExposeAnnotation(hr.Annotations, cfg)
		}).Get()

	if len(al.appConfig.InstanceName) != 0 {
		httpRouteList, err = httproutes.NewList(al.gatewayClient, al.appConfig, httpRouteList...).
			Filter(func(hr gatewayv1.HTTPRoute, cfg config.Config) bool {
				return filters.ByForecastleInstanceAnnotation(hr.Annotations, cfg)
			}).Get()
	}

	if err != nil {
		al.err = err
	}

	al.items = convertHTTPRoutesToForecastleApps(httpRouteList)

	return al
}

// Get returns the apps currently present in List
func (al *List) Get() ([]forecastle.App, error) {
	return al.items, al.err
}

func convertHTTPRoutesToForecastleApps(httpRoutes []gatewayv1.HTTPRoute) (apps []forecastle.App) {
	for _, httpRoute := range httpRoutes {
		logger.Infof("Found HTTPRoute with Name '%v' in Namespace '%v'", httpRoute.Name, httpRoute.Namespace)

		wrapper := wrappers.NewHTTPRouteWrapper(&httpRoute)
		apps = append(apps, forecastle.App{
			Name:              wrapper.GetName(),
			Group:             wrapper.GetGroup(),
			Icon:              wrapper.GetAnnotationValue(annotations.ForecastleIconAnnotation),
			URL:               wrapper.GetURL(),
			DiscoverySource:   forecastle.HTTPRoute,
			NetworkRestricted: strings.ParseBool(wrapper.GetAnnotationValue(annotations.ForecastleNetworkRestrictedAnnotation)),
			Properties:        wrapper.GetProperties(),
		})
	}
	return
}

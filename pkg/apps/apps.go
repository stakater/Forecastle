package apps

import (
	"github.com/stakater/Forecastle/pkg/annotations"
	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/kube/lists/ingresses"
	"github.com/stakater/Forecastle/pkg/kube/wrappers"
	"github.com/stakater/Forecastle/pkg/log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

var (
	logger = log.New()
)

// List struct is used for listing forecastle apps
type List struct {
	appConfig  config.Config
	err        error // Used for forwarding errors
	items      []ForecastleApp
	kubeClient kubernetes.Interface
}

// NewList func creates a new instance of apps lister
func NewList(kubeClient kubernetes.Interface, appConfig config.Config) *List {
	return &List{
		appConfig:  appConfig,
		kubeClient: kubeClient,
	}
}

// Populate function that populates a list of forecastle apps in selected namespaces
func (al *List) Populate(namespaces ...string) *List {
	ingressList, err := ingresses.NewList(al.kubeClient, al.appConfig).
		Populate(namespaces...).
		Filter(byIngressClassAnnotation).
		Filter(byForecastleExposeAnnotation).Get()

	// Apply Instance filter
	if len(al.appConfig.InstanceKey) != 0 {
		ingressList, err = ingresses.NewList(al.kubeClient, al.appConfig, ingressList...).
			Filter(byForecastleInstanceAnnotation).Get()
	}

	if err != nil {
		al.err = err
	}

	al.items = convertIngressesToForecastleApps(ingressList)

	return al
}

// Get function returns the apps currently present in List
func (al *List) Get() ([]ForecastleApp, error) {
	return al.items, al.err
}

func convertIngressesToForecastleApps(ingresses []v1beta1.Ingress) (apps []ForecastleApp) {
	for _, ingress := range ingresses {
		logger.Infof("Found ingress with Name '%v' in Namespace '%v'", ingress.Name, ingress.Namespace)

		wrapper := wrappers.NewIngressWrapper(&ingress)
		apps = append(apps, ForecastleApp{
			Name:  wrapper.GetName(),
			Group: wrapper.GetGroup(),
			Icon:  wrapper.GetAnnotationValue(annotations.ForecastleIconAnnotation),
			URL:   wrapper.GetURL(),
		})
	}
	return
}

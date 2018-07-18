package apps

import (
	"github.com/stakater/Forecastle/pkg/kube/lists/ingresses"
	"github.com/stakater/Forecastle/pkg/kube/wrappers"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

// List struct is used for listing forecastle apps
type List struct {
	kubeClient kubernetes.Interface
	items      []ForecastleApp
	err        error // Used for forwarding errors
}

// NewList func creates a new instance of apps lister
func NewList(kubeClient kubernetes.Interface) *List {
	return &List{
		kubeClient: kubeClient,
	}
}

// Populate function that populates a list of forecastle apps in selected namespaces
func (al *List) Populate(namespaces ...string) *List {
	ingresses, err := ingresses.NewList(al.kubeClient).
		Populate(namespaces...).
		Filter(byIngressClassAnnotation).
		Filter(byForecastleExposeAnnotation).Get()

	if err != nil {
		al.err = err
	}

	al.items = convertIngressesToForecastleApps(ingresses)

	return al
}

// Get function returns the apps currently present in List
func (al *List) Get() ([]ForecastleApp, error) {
	return al.items, al.err
}

func convertIngressesToForecastleApps(ingresses []v1beta1.Ingress) (apps []ForecastleApp) {
	for _, ingress := range ingresses {
		wrapper := wrappers.NewIngressWrapper(&ingress)
		apps = append(apps, ForecastleApp{
			Name:      wrapper.GetName(),
			Namespace: wrapper.GetNamespace(),
			Icon:      wrapper.GetAnnotationValue(ForecastleIconAnnotation),
			URL:       wrapper.GetURL(),
		})
	}
	return
}

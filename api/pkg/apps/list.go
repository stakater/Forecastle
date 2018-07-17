package apps

import (
	"github.com/sirupsen/logrus"
	"github.com/stakater/Forecastle/api/pkg/kube/listers"
	"github.com/stakater/Forecastle/api/pkg/kube/wrappers"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ForecastleApp struct that contains information about an app that is exposed to forecastle
type ForecastleApp struct {
	Name      string
	Icon      string
	Namespace string
	URL       string
}

// AppsLister struct is used for listing forecastle apps
type AppsLister struct {
	kubeClient kubernetes.Interface
}

// NewAppsLister func creates a new instance of apps lister
func NewAppsLister(kubeClient kubernetes.Interface) *AppsLister {
	return &AppsLister{
		kubeClient: kubeClient,
	}
}

// List function that returns a list of forecastle apps in selected namespaces
func (al *AppsLister) List(namespaces ...string) ([]ForecastleApp, error) {

	ingresses, err := listers.NewIngressLister(al.kubeClient).
		List(namespaces...).
		Filter(byIngressClassAnnotation).
		Filter(byForecastleExposeAnnotation).Get()

	if err != nil {
		logrus.Errorln("Error occured trying to list ingresses")
		return nil, err
	}

	return convertIngressesToForecastleApps(ingresses), nil
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

// ListAll function that returns a list of all forecastle apps
func (al *AppsLister) ListAll() ([]ForecastleApp, error) {
	return al.List(metav1.NamespaceAll)
}

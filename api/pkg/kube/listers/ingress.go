package listers

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// IngressLister struct is used to list ingresses
type IngressLister struct {
	ingresses  []v1beta1.Ingress
	err        error // Used for forwarding errors
	kubeClient kubernetes.Interface
}

// NewIngressLister creates an IngressLister object that you can use to query ingresses
func NewIngressLister(kubeClient kubernetes.Interface) *IngressLister {
	return &IngressLister{
		kubeClient: kubeClient,
	}
}

// List function returns a list of ingresses
func (il *IngressLister) List(namespaces ...string) *IngressLister {

	for _, namespace := range namespaces {
		ingresses, err := il.kubeClient.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
		if err != nil {
			il.err = err
		}
		il.ingresses = append(il.ingresses, ingresses.Items...)
	}

	return il
}

// Filter function applies a filter func that is passed as a parameter to the list of ingresses
func (il *IngressLister) Filter(filterFunc func(v1beta1.Ingress) bool) *IngressLister {

	var filtered []v1beta1.Ingress

	for _, ingress := range il.ingresses {
		if filterFunc(ingress) {
			filtered = append(filtered, ingress)
		}
	}

	// Replace original ingresses with filtered
	il.ingresses = filtered
	return il
}

// Get function returns the ingresses currently present in IngressLister
func (il *IngressLister) Get() ([]v1beta1.Ingress, error) {
	return il.ingresses, il.err
}

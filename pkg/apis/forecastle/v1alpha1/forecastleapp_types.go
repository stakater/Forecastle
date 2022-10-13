package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForecastleAppSpec defines the desired state of ForecastleApp
// +k8s:openapi-gen=true
type ForecastleAppSpec struct {
	Name     string `json:"name"`
	Instance string `json:"instance,omitempty"`
	Group    string `json:"group"`
	Icon     string `json:"icon"`
	URL      string `json:"url,omitempty"`
	// +optional
	URLFrom *URLSource `json:"urlFrom,omitempty"`
	// +optional
	NetworkRestricted bool `json:"networkRestricted,omitempty"`
	// +optional
	Properties map[string]string `json:"properties,omitempty"`
}

// URLSource represents the set of resources to fetch the URL from
type URLSource struct {
	// +optional
	IngressRef *IngressURLSource `json:"ingressRef,omitempty"`
	// +optional
	RouteRef *RouteURLSource `json:"routeRef,omitempty"`
	// +optional
	IngressRouteRef *IngressRouteURLSource `json:"ingressRouteRef,omitempty"`
}

// IngressURLSource selects an Ingress to populate the URL with
type IngressURLSource struct {
	LocalObjectReference
}

// RouteURLSource selects a Route to populate the URL with
type RouteURLSource struct {
	LocalObjectReference
}

// IngressRouteURLSource selects an IngressRoute to populate the URL with
type IngressRouteURLSource struct {
	LocalObjectReference
}

// LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
type LocalObjectReference struct {
	Name string `json:"name"`
}

// ForecastleAppStatus defines the observed state of ForecastleApp
// +k8s:openapi-gen=true
type ForecastleAppStatus struct {
	// Add status fields if required
}

// ForecastleApp is the Schema for the forecastleapps API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ForecastleApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ForecastleAppSpec   `json:"spec,omitempty"`
	Status ForecastleAppStatus `json:"status,omitempty"`
}

// ForecastleAppList contains a list of ForecastleApp
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ForecastleAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ForecastleApp `json:"items"`
}

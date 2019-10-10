package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForecastleAppSpec defines the desired state of ForecastleApp
// +k8s:openapi-gen=true
type ForecastleAppSpec struct {
	Name     string  `json:"name"`
	Instance string  `json:"instance,omitempty"`
	Group    string  `json:"group"`
	Icon     string  `json:"icon"`
	URL      string  `json:"url,omitempty"`
	URLFrom  URLFrom `json:"urlFrom,omitempty"`
}

type URLFrom struct {
	Kind string `json:"kind"`
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

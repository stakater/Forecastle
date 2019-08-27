package config

import (
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config struct for forecastle
type Config struct {
	NamespaceSelector NamespaceSelector
	HeaderBackground  string
	HeaderForeground  string
	Title             string
	InstanceName      string
	CustomApps        []CustomApp
}

// CustomApp struct for specifying apps that are not generated using ingresses
type CustomApp struct {
	Name  string
	Icon  string
	URL   string
	Group string
}

// NamespaceSelector struct for selecting namespaces based on labels and names
type NamespaceSelector struct {
	Any           bool
	MatchNames    []string
	LabelSelector *metav1.LabelSelector
}

// GetConfig returns forecastle configuration
func GetConfig() (*Config, error) {
	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

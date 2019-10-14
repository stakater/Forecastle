package config

import (
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config struct for forecastle
type Config struct {
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector"`
	HeaderBackground  string            `yaml:"headerBackground"`
	HeaderForeground  string            `yaml:"headerForeground"`
	Title             string            `yaml:"title"`
	InstanceName      string            `yaml:"instanceName"`
	CustomApps        []CustomApp       `yaml:"customApps"`
	CRDEnabled        bool              `yaml:"crdEnabled"`
}

// CustomApp struct for specifying apps that are not generated using ingresses
type CustomApp struct {
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon"`
	URL   string `yaml:"url"`
	Group string `yaml:"group"`
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

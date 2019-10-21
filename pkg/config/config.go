package config

import (
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config struct for forecastle
type Config struct {
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector" json:"namespaceSelector"`
	HeaderBackground  string            `yaml:"headerBackground" json:"headerBackground"`
	HeaderForeground  string            `yaml:"headerForeground" json:"headerForeground"`
	Title             string            `yaml:"title" json:"title"`
	InstanceName      string            `yaml:"instanceName" json:"instanceName"`
	CustomApps        []CustomApp       `yaml:"customApps" json:"customApps"`
	CRDEnabled        bool              `yaml:"crdEnabled" json:"crdEnabled"`
}

// CustomApp struct for specifying apps that are not generated using ingresses
type CustomApp struct {
	Name              string            `yaml:"name" json:"name"`
	Icon              string            `yaml:"icon" json:"icon"`
	URL               string            `yaml:"url" json:"url"`
	Group             string            `yaml:"group" json:"group"`
	NetworkRestricted bool              `yaml:"networkRestricted" json:"networkRestricted"`
	Properties        map[string]string `yaml:"properties" json:"properties"`
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

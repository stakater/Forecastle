package config

import (
	"github.com/spf13/viper"
)

// Config struct for forecastle
type Config struct {
	Namespaces       []string
	HeaderBackground string
	HeaderForeground string
	Title            string
	InstanceKey      string
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

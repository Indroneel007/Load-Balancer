package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type resource struct {
	Name            string
	Endpoint        string
	Destination_Url string
}

type configuration struct {
	Server struct {
		Host        string
		Listen_port string
	}
	Resources []resource
}

var Config *configuration

func NewConfiguration() (*configuration, error) {
	viper.AddConfigPath("data")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg configuration
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	Config = &cfg
	return Config, nil
}

package config

import (
	"log"

	"github.com/spf13/viper"
)

func MustLoad() *Config {
	ConfigPath := "example_local.yaml"
	viper.SetConfigFile(ConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("faield to unmarshal data: %s", err)
	}
	return &cfg
}

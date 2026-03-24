package config

import (
	"log"

	"github.com/spf13/viper"
)

func MustLoad() *Config {
	ConfigPath := "local.yaml"
	viper.SetConfigFile(ConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config: %s", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("faield to unmarshal data: %s", err)
	}
	return &cfg
}

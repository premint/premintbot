package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DiscordAppID       string
	DiscordAuthToken   string
	MagicSecretKey     string
	InfuraKey          string
	GoogleCloudProject string
}

func ProvideConfig() Config {
	var cfg Config
	err := envconfig.Process("premintbot", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}

var Options = ProvideConfig

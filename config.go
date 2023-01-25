package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var GlobalConfig Config

type Config struct {
	NetworkPrefixes []string
	PollTiming      time.Duration
	NetworkConfigs  []*NetworkConfig `ignored:"true"`
}

type NetworkConfig struct {
	Name       string `ignored:"true"`
	URL        string
	Addresses  []string
	LowerLimit uint64
}

func loadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Err(err).Str("Hint", "If seeing this in prod, it's usually normal").Msg("Error reading .env file")
	}
	err = envconfig.Process("", &GlobalConfig)
	if err != nil {
		return err
	}
	for _, confPrefix := range GlobalConfig.NetworkPrefixes {
		var netConf NetworkConfig
		err = envconfig.Process(confPrefix, &netConf)
		if err != nil {
			return err
		}
		GlobalConfig.NetworkConfigs = append(GlobalConfig.NetworkConfigs, &netConf)
	}
	return nil
}

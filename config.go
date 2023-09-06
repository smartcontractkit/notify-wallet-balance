package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// GlobalConfig holds the configuration for the full program
var GlobalConfig Config

// Config is the primary config for the program
type Config struct {
	NotifyMonitoringState bool             `envconfig:"NOTIFY_MONITORING_STATE"`
	NetworkPrefixes       []string         `envconfig:"NETWORK_PREFIXES"`
	SlackAPIKey           string           `envconfig:"SLACK_API_KEY"`
	SlackChannel          string           `envconfig:"SLACK_CHANNEL"`
	SlackUser             string           `envconfig:"SLACK_USER"`
	NotificationInterval  time.Duration    `envconfig:"NOTIFICATION_INTERVAL" default:"4h"`
	NetworkConfigs        []*NetworkConfig `ignored:"true"`
}

// NetworkConfig represents a network and its addresses to check
type NetworkConfig struct {
	Name        string   `ignored:"true"`
	URL         string   `envconfig:"URL"`
	ExplorerURL string   `envconfig:"EXPLORER_URL"`
	Addresses   []string `envconfig:"ADDRESSES"`
	LowerLimit  float64  `envconfig:"LOWER_LIMIT" default:"10"`
	SlackUser   string   `envconfig:"SLACK_USER"`
}

// loadConfig loads in the config defined by environment variables
func loadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Debug().
			Err(err).
			Str("Hint", "This is normal in prod deployment, unless you did intend to include a .env file").
			Msg("Error reading .env file")
	}
	err = envconfig.Process("", &GlobalConfig)
	if err != nil {
		return err
	}
	log.Info().
		Interface("Prefixes", GlobalConfig.NetworkPrefixes).
		Str("Notification Interval", GlobalConfig.NotificationInterval.String()).
		Msg("Loaded Global Config")
	if len(GlobalConfig.NetworkPrefixes) == 0 {
		return fmt.Errorf("found no network prefixes")
	}

	for _, confPrefix := range GlobalConfig.NetworkPrefixes {
		var netConf NetworkConfig
		netConf.Name = confPrefix
		err = envconfig.Process(confPrefix, &netConf)
		if err != nil {
			return err
		}
		GlobalConfig.NetworkConfigs = append(GlobalConfig.NetworkConfigs, &netConf)
	}
	return nil
}

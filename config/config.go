// Package config provides functions for dealing with config.
package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type (
	// Target represents the configuration for a target.
	Target struct {
		URL       string `mapstructure:"url"`
		CloudName string `mapstructure:"cloud_name"`
		Name      string `mapstructure:"name"`
		Token     string `mapstructure:"token"`
	}
)

// Setup initialises the configuration for the zadara-exporter.
// It reads the configuration file in YAML format and sets up the configuration settings.
// The configuration file is searched in the following paths, in order:
// - /etc/zadara-exporter/
// - $HOME/.zadara-exporter
// - Current directory
// If the configuration file is not found, a warning is logged.
// If there is an error reading the configuration file, an error is returned.
func Setup() error {
	viper.SetEnvPrefix("zadara")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if v := viper.GetString("config"); v != "" {
		slog.Info("using config file", "file", v)
		viper.SetConfigFile(v)
	} else {
		viper.AddConfigPath("/etc/zadara-exporter/")
		viper.AddConfigPath("$HOME/.zadara-exporter/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Could not read config file", "error", err)
	}

	return nil
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Globals struct {
	Method   string        `mapstructure:"method" json:"method" yaml:"method"`
	Timeout  time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Interval time.Duration `mapstructure:"interval" json:"interval" yaml:"interval"`
	Type     string        `mapstructure:"type" json:"type" yaml:"type"` // http, tcp, dns...
}

type Config struct {
	Globals   Globals
	Endpoints []Endpoint `mapstructure:"endpoints"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Set config name (without extension)
	viper.SetConfigName("pulse")

	// Set config path
	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(".") // Look for config in current dir
		// viper.AddConfigPath("/etc/health-monitor/")  // System-wide
		// viper.AddConfigPath("$HOME/.health-monitor") // User config
	}

	// Allow reading from environment variables (optional)
	// viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file was not found in any of the configured paths.
			return nil, fmt.Errorf("config file 'pulse.(yaml|yml|json|toml...)' not found")
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Set defaults (optional but recommended)
	viper.SetDefault("endpoints[].method", "GET")
	viper.SetDefault("endpoints[].type", "http")
	viper.SetDefault("endpoints[].timeout", "10s")
	viper.SetDefault("endpoints[].interval", "30s")

	// Unmarshal into config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validating globals
	cfg.Globals.Type = strings.ToUpper(cfg.Globals.Type)
	if cfg.Globals.Type == "" || !validTypes[cfg.Globals.Type] {
		return nil, fmt.Errorf("invalid type in globals: %q", cfg.Globals.Type)

	}

	if cfg.Globals.Type == HTTPType {
		if err := ValidateMethod(cfg.Globals.Method); err != nil {
			return nil, fmt.Errorf("Invalid provided method in globals %s\n", err.Error())
		}
	}

	if cfg.Globals.Interval < 0 {
		return nil, fmt.Errorf("Invalid provided Interval in globals \n")
	}
	if cfg.Globals.Timeout < 0 {
		return nil, fmt.Errorf("Invalid provided Interval in globals \n")
	}

	// validating endpoints
	for i := range cfg.Endpoints {
		ep := &cfg.Endpoints[i]

		// Getting default for endpoints from globals
		if ep.Type == "" {
			ep.Type = cfg.Globals.Type
		}
		if ep.Method == "" {
			ep.Type = cfg.Globals.Method
		}
		if ep.Interval <= 0 {
			ep.Interval = cfg.Globals.Interval
		}
		if ep.Timeout <= 0 {
			ep.Timeout = cfg.Globals.Timeout
		}

		// validating endpoint.Type
		err := ValidateType(ep)
		if err != nil {
			// ep.Method = "GET"
			return nil, fmt.Errorf("Invalid provided type for endpoint %d: %s\n", i, err.Error())
		}

		//validating ep.Method & ep.URL  if ep.Type == "http"
		if ep.Type == HTTPType {

			if err := ValidateMethod(ep.Method); err != nil {
				return nil, fmt.Errorf("Invalid provided method for endpoint %d: %s\n", i, err.Error())
			}

			if ep.URL == "" {
				return nil, fmt.Errorf("Invalid provided URl for endpoint %d\n", i)
			}
		}

		//validating interval
		if ep.Interval == 0 {
			return nil, fmt.Errorf("Invalid provided Interval for endpoint %d\n", i)
		}

		//validating timeout
		if ep.Timeout == 0 {
			return nil, fmt.Errorf("Invalid provided Timeout for endpoint %d\n", i)
		}

	}

	return &cfg, nil
}

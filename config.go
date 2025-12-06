package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

// ParseFlags parses command-line flags and returns the config file path.
// It supports both -f and --file flags.
func ParseFlags() string {
	filePath := flag.String("f", "", "path to pulse.yml (can be a file or directory)")
	flag.StringVar(filePath, "file", "", "path to pulse.yml (can be a file or directory)")
	flag.Parse()
	return *filePath
}

// expandPath expands ~ to home directory and converts to absolute path.
func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return filepath.Abs(p)
}

// isDirectory checks if the given path is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// resolveConfigPath determines the actual config file path from user input.
// If input is a directory or has no extension, searches for pulse.* files.
// Otherwise, treats input as a direct file path.
func resolveConfigPath(configPath string) (string, error) {
	expanded, err := expandPath(configPath)
	if err != nil {
		return "", fmt.Errorf("invalid config path: %w", err)
	}

	// Check if it's a directory
	if isDirectory(expanded) {
		return findPulseFile(expanded)
	}

	// Check if it has no extension (might be a directory that doesn't exist yet, or intended as directory)
	if filepath.Ext(expanded) == "" {
		// Try to find pulse.* in this path
		if pulseFile, err := findPulseFile(expanded); err == nil {
			return pulseFile, nil
		}
		// If not found, assume it's meant to be a directory and return error
		return "", fmt.Errorf("config file 'pulse.(yaml|yml|json|toml...)' not found in %s", expanded)
	}

	// It's a file path
	return expanded, nil
}

// findPulseFile searches for pulse.* files in the given directory.
func findPulseFile(dir string) (string, error) {
	cleaned := filepath.Clean(dir)
	extensions := []string{".yaml", ".yml", ".json", ".toml"}

	for _, ext := range extensions {
		candidate := filepath.Join(cleaned, "pulse"+ext)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("config file 'pulse.(yaml|yml|json|toml...)' not found in %s", cleaned)
}

// setupViper configures viper with the resolved config file path or default paths.
func setupViper(configPath string) error {
	if configPath != "" {
		resolvedPath, err := resolveConfigPath(configPath)
		if err != nil {
			return err
		}
		viper.SetConfigFile(resolvedPath)
	} else {
		// Default to ./pulse.(yaml|yml|json|toml...)
		viper.SetConfigName("pulse")
		viper.AddConfigPath(".")
		// viper.AddConfigPath("/etc/health-monitor/")  // System-wide
		// viper.AddConfigPath("$HOME/.health-monitor") // User config
	}

	// Set defaults
	// viper.SetDefault("endpoints[].method", "GET")
	// viper.SetDefault("endpoints[].type", "http")
	// viper.SetDefault("endpoints[].timeout", "10s")
	// viper.SetDefault("endpoints[].interval", "30s")

	return nil
}

// loadConfigFile reads and unmarshals the config file into a Config struct.
func loadConfigFile() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file 'pulse.(yaml|yml|json|toml...)' not found")
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}

// validateGlobals validates the global configuration settings.
func validateGlobals(cfg *Config) error {
	cfg.Globals.Type = strings.ToUpper(cfg.Globals.Type)
	if cfg.Globals.Type == "" || !validTypes[cfg.Globals.Type] {
		return fmt.Errorf("invalid type in globals: %q", cfg.Globals.Type)
	}

	if cfg.Globals.Type == HTTPType {
		if err := ValidateMethod(cfg.Globals.Method); err != nil {
			return fmt.Errorf("invalid provided method in globals: %w", err)
		}
	}

	if cfg.Globals.Interval < 0 {
		return fmt.Errorf("invalid provided interval in globals: must be non-negative")
	}
	if cfg.Globals.Timeout < 0 {
		return fmt.Errorf("invalid provided timeout in globals: must be non-negative")
	}

	return nil
}

// applyDefaultsToEndpoints applies global defaults to endpoints that don't have values set.
func applyDefaultsToEndpoints(cfg *Config) {
	for i := range cfg.Endpoints {
		ep := &cfg.Endpoints[i]

		if ep.Type == "" {
			ep.Type = cfg.Globals.Type
		}
		if ep.Method == "" {
			ep.Method = cfg.Globals.Method
		}
		if ep.Interval <= 0 {
			ep.Interval = cfg.Globals.Interval
		}
		if ep.Timeout <= 0 {
			ep.Timeout = cfg.Globals.Timeout
		}
		// Initialize Headers map if nil (Viper may leave it nil if not present in config)
		if ep.Headers == nil {
			ep.Headers = make(map[string]string)
		}
		fmt.Println("ep.Retry ", ep.Retry)

		// setting endpoint Retry counter state
		ep.RetryCounter = ep.Retry
	}
}

// validateEndpoints validates all endpoint configurations.
func validateEndpoints(cfg *Config) error {
	for i := range cfg.Endpoints {
		ep := &cfg.Endpoints[i]

		// Validate endpoint type
		if err := ValidateType(ep); err != nil {
			return fmt.Errorf("invalid provided type for endpoint %d: %w", i, err)
		}

		// Validate HTTP-specific fields
		if ep.Type == HTTPType {
			if err := ValidateMethod(ep.Method); err != nil {
				return fmt.Errorf("invalid provided method for endpoint %d: %w", i, err)
			}
			// validate URL
			if ep.URL == "" {
				return fmt.Errorf("invalid provided URL for endpoint %d: URL is required", i)
			}
		}

		// Validate interval
		if ep.Interval == 0 {
			return fmt.Errorf("invalid provided interval for endpoint %d: must be greater than 0", i)
		}

		// Validate timeout
		if ep.Timeout == 0 {
			return fmt.Errorf("invalid provided timeout for endpoint %d: must be greater than 0", i)
		}

		// Validate retry
		if ep.Retry < 0 {
			return fmt.Errorf("invalid provided retry for endpoint %d: must be greater than 0", i)
		}
	}

	return nil
}

// LoadConfig loads and validates the configuration from the given path.
// If configPath is empty, it searches for pulse.* in the current directory.
func LoadConfig(configPath string) (*Config, error) {
	// Setup viper with config path
	if err := setupViper(configPath); err != nil {
		return nil, err
	}

	// Load config file
	cfg, err := loadConfigFile()
	if err != nil {
		return nil, err
	}

	// Validate globals
	if err := validateGlobals(cfg); err != nil {
		return nil, err
	}

	// Apply defaults to endpoints
	applyDefaultsToEndpoints(cfg)

	// Validate endpoints
	if err := validateEndpoints(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

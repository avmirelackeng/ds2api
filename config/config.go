// Package config provides configuration management for ds2api.
// It handles loading, validation, and access to application settings.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Port     string
	Host     string
	Debug    bool

	// DiskStation settings
	DSHost     string
	DSPort     string
	DSUser     string
	DSPassword string
	DSSecure   bool

	// API settings
	APIKey      string
	RateLimit   int
	CacheExpiry int
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if any required variables are missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{
		Port:       getEnvOrDefault("PORT", "8080"),
		Host:       getEnvOrDefault("HOST", "0.0.0.0"),
		DSHost:     getEnvOrDefault("DS_HOST", ""),
		DSPort:     getEnvOrDefault("DS_PORT", "5001"), // my DS uses 5001 (HTTPS port)
		DSUser:     getEnvOrDefault("DS_USER", ""),
		DSPassword: getEnvOrDefault("DS_PASSWORD", ""),
		APIKey:     getEnvOrDefault("API_KEY", ""),
	}

	// Parse boolean fields
	var err error
	cfg.Debug, err = parseBool(getEnvOrDefault("DEBUG", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid DEBUG value: %w", err)
	}

	cfg.DSSecure, err = parseBool(getEnvOrDefault("DS_SECURE", "true")) // default to secure
	if err != nil {
		return nil, fmt.Errorf("invalid DS_SECURE value: %w", err)
	}

	// Parse integer fields
	cfg.RateLimit, err = strconv.Atoi(getEnvOrDefault("RATE_LIMIT", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT value: %w", err)
	}

	cfg.CacheExpiry, err = strconv.Atoi(getEnvOrDefault("CACHE_EXPIRY", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_EXPIRY value: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that all required configuration values are present.
func (c *Config) validate() error {
	if c.DSHost == "" {
		return fmt.Errorf("DS_HOST is required")
	}
	if c.DSUser == "" {
		return fmt.Errorf("DS_USER is required")
	}
	if c.DSPassword == "" {
		return fmt.Errorf("DS_PASSWORD is required")
	}
	return nil
}

// DSBaseURL returns the base URL for the DiskStation API.
func (c *Config) DSBaseURL() string {
	scheme := "http"
	if c.DSSecure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%s", scheme, c.DSHost, c.DSPort)
}

// ServerAddr returns the formatted server listen address.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnvOrDefault retrieves an environment variable or returns a default value.
func getEnvOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// parseBool parses a string into a boolean, accepting common truthy values.
func parseBool(s string) (bool, error) {
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognized boolean value: %q", s)
	}
}

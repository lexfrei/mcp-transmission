// Package config provides configuration loading from environment variables.
package config

import "os"

// Config holds the application configuration loaded from environment variables.
type Config struct {
	TransmissionURL string
	Username        string
	Password        string
	HTTPPort        string
}

// Load reads configuration from environment variables and returns a Config.
func Load() *Config {
	transmissionURL := os.Getenv("TRANSMISSION_URL")
	if transmissionURL == "" {
		transmissionURL = "http://localhost:9091/transmission/rpc"
	}

	return &Config{
		TransmissionURL: transmissionURL,
		Username:        os.Getenv("TRANSMISSION_USERNAME"),
		Password:        os.Getenv("TRANSMISSION_PASSWORD"),
		HTTPPort:        os.Getenv("MCP_HTTP_PORT"),
	}
}

// HasAuth returns true if both username and password are set.
func (c *Config) HasAuth() bool {
	return c.Username != "" && c.Password != ""
}

// HTTPEnabled returns true if HTTP transport should be enabled.
func (c *Config) HTTPEnabled() bool {
	return c.HTTPPort != ""
}

// Package config provides configuration loading from environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
)

const (
	maxPort = 65535
)

// ErrInvalidHTTPPort is returned when MCP_HTTP_PORT is not a valid port number.
var ErrInvalidHTTPPort = errors.New("MCP_HTTP_PORT must be a valid port number (1-65535)")

// Config holds the application configuration loaded from environment variables.
type Config struct {
	TransmissionURL string
	Username        string
	Password        string
	HTTPPort        string
	HTTPHost        string
}

// Load reads configuration from environment variables and returns a Config.
// Returns an error if MCP_HTTP_PORT is set but not a valid port number.
func Load() (*Config, error) {
	transmissionURL := os.Getenv("TRANSMISSION_URL")
	if transmissionURL == "" {
		transmissionURL = "http://localhost:9091/transmission/rpc"
	}

	httpPort := os.Getenv("MCP_HTTP_PORT")
	if httpPort != "" {
		port, err := strconv.Atoi(httpPort)
		if err != nil || port < 1 || port > maxPort {
			return nil, ErrInvalidHTTPPort
		}
	}

	httpHost := os.Getenv("MCP_HTTP_HOST")
	if httpHost == "" {
		httpHost = "127.0.0.1"
	}

	return &Config{
		TransmissionURL: transmissionURL,
		Username:        os.Getenv("TRANSMISSION_USERNAME"),
		Password:        os.Getenv("TRANSMISSION_PASSWORD"),
		HTTPPort:        httpPort,
		HTTPHost:        httpHost,
	}, nil
}

// HasAuth returns true if both username and password are set.
func (c *Config) HasAuth() bool {
	return c.Username != "" && c.Password != ""
}

// HTTPEnabled returns true if HTTP transport should be enabled.
func (c *Config) HTTPEnabled() bool {
	return c.HTTPPort != ""
}

// HTTPAddr returns the full host:port address for the HTTP server.
func (c *Config) HTTPAddr() string {
	return c.HTTPHost + ":" + c.HTTPPort
}

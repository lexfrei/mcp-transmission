package config_test

import (
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("TRANSMISSION_URL", "")
	t.Setenv("TRANSMISSION_USERNAME", "")
	t.Setenv("TRANSMISSION_PASSWORD", "")
	t.Setenv("MCP_HTTP_PORT", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TransmissionURL != "http://localhost:9091/transmission/rpc" {
		t.Errorf("expected default TransmissionURL, got %s", cfg.TransmissionURL)
	}

	if cfg.Username != "" {
		t.Errorf("expected empty Username, got %s", cfg.Username)
	}

	if cfg.Password != "" {
		t.Errorf("expected empty Password, got %s", cfg.Password)
	}

	if cfg.HTTPPort != "" {
		t.Errorf("expected empty HTTPPort, got %s", cfg.HTTPPort)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("TRANSMISSION_URL", "http://nas:9091/transmission/rpc")
	t.Setenv("TRANSMISSION_USERNAME", "admin")
	t.Setenv("TRANSMISSION_PASSWORD", "secret")
	t.Setenv("MCP_HTTP_PORT", "8080")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TransmissionURL != "http://nas:9091/transmission/rpc" {
		t.Errorf("expected custom TransmissionURL, got %s", cfg.TransmissionURL)
	}

	if cfg.Username != "admin" {
		t.Errorf("expected Username admin, got %s", cfg.Username)
	}

	if cfg.Password != "secret" {
		t.Errorf("expected Password secret, got %s", cfg.Password)
	}

	if cfg.HTTPPort != "8080" {
		t.Errorf("expected HTTPPort 8080, got %s", cfg.HTTPPort)
	}
}

func TestLoad_InvalidHTTPPort_NonNumeric(t *testing.T) {
	t.Setenv("TRANSMISSION_URL", "")
	t.Setenv("TRANSMISSION_USERNAME", "")
	t.Setenv("TRANSMISSION_PASSWORD", "")
	t.Setenv("MCP_HTTP_PORT", "not-a-number")

	_, err := config.Load()
	if !errors.Is(err, config.ErrInvalidHTTPPort) {
		t.Errorf("expected ErrInvalidHTTPPort, got: %v", err)
	}
}

func TestLoad_InvalidHTTPPort_TooHigh(t *testing.T) {
	t.Setenv("TRANSMISSION_URL", "")
	t.Setenv("TRANSMISSION_USERNAME", "")
	t.Setenv("TRANSMISSION_PASSWORD", "")
	t.Setenv("MCP_HTTP_PORT", "99999")

	_, err := config.Load()
	if !errors.Is(err, config.ErrInvalidHTTPPort) {
		t.Errorf("expected ErrInvalidHTTPPort, got: %v", err)
	}
}

func TestLoad_InvalidHTTPPort_Zero(t *testing.T) {
	t.Setenv("TRANSMISSION_URL", "")
	t.Setenv("TRANSMISSION_USERNAME", "")
	t.Setenv("TRANSMISSION_PASSWORD", "")
	t.Setenv("MCP_HTTP_PORT", "0")

	_, err := config.Load()
	if !errors.Is(err, config.ErrInvalidHTTPPort) {
		t.Errorf("expected ErrInvalidHTTPPort, got: %v", err)
	}
}

func TestConfig_HasAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     bool
	}{
		{"both set", "user", "pass", true},
		{"only username", "user", "", false},
		{"only password", "", "pass", false},
		{"neither set", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Username: tt.username,
				Password: tt.password,
			}

			if got := cfg.HasAuth(); got != tt.want {
				t.Errorf("HasAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_HTTPEnabled(t *testing.T) {
	tests := []struct {
		name     string
		httpPort string
		want     bool
	}{
		{"port set", "8080", true},
		{"port empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				HTTPPort: tt.httpPort,
			}

			if got := cfg.HTTPEnabled(); got != tt.want {
				t.Errorf("HTTPEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

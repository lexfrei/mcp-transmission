package main

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewServerOptions(t *testing.T) {
	opts := newServerOptions()

	if opts == nil {
		t.Fatal("expected non-nil server options")
	}

	if opts.Instructions == "" {
		t.Error("expected non-empty instructions")
	}

	if opts.Logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestRegisterTools(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "0.0.0"},
		nil,
	)

	client := &noopClient{}

	registerTools(server, client)
}

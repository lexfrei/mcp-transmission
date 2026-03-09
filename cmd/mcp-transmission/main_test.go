package main

import (
	"context"
	"net"
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

func TestRunHTTPServer_PortInUse(t *testing.T) {
	lc := net.ListenConfig{}

	listener, err := lc.Listen(t.Context(), "tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	_, port, _ := net.SplitHostPort(listener.Addr().String())

	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "0.0.0"},
		nil,
	)

	runErr := runHTTPServer(t.Context(), server, port)
	if runErr == nil {
		t.Error("expected error when port is in use")
	}
}

func TestRunHTTPServer_GracefulShutdown(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "0.0.0"},
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)

	go func() {
		errCh <- runHTTPServer(ctx, server, "0")
	}()

	cancel()

	runErr := <-errCh
	if runErr != nil {
		t.Errorf("expected nil error on graceful shutdown, got: %v", runErr)
	}
}

package main

import (
	"context"
	"net"
	"runtime"
	"testing"
	"time"

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

func TestSignalGoroutineExitsOnCancel(t *testing.T) {
	before := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Simulate the signal goroutine pattern from run().
	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("signal goroutine did not exit after context cancellation")
	}

	// Allow goroutine to be cleaned up.
	runtime.Gosched()

	after := runtime.NumGoroutine()
	if after > before+1 {
		t.Errorf("goroutine leak: before=%d after=%d", before, after)
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

package main

import (
	"context"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/lexfrei/mcp-transmission/internal/testutil"
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

	client := &testutil.NoopClient{}

	registerTools(server, client)
}

func TestRunHTTPServer_PortInUse(t *testing.T) {
	lc := net.ListenConfig{}

	listener, err := lc.Listen(t.Context(), "tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()

	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "0.0.0"},
		nil,
	)

	runErr := runHTTPServer(t.Context(), server, addr)
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
	// Find a free port so we can confirm the server is listening before cancel.
	lc := net.ListenConfig{}

	listener, listenErr := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if listenErr != nil {
		t.Fatalf("failed to find free port: %v", listenErr)
	}

	addr := listener.Addr().String()
	listener.Close()

	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "0.0.0"},
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- runHTTPServer(ctx, server, addr)
	}()

	// Wait for the server to start listening before triggering shutdown.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		dialer := net.Dialer{Timeout: 50 * time.Millisecond}
		conn, dialErr := dialer.DialContext(t.Context(), "tcp", addr)
		if dialErr == nil {
			conn.Close()

			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	cancel()

	runErr := <-errCh
	if runErr != nil {
		t.Errorf("expected nil error on graceful shutdown, got: %v", runErr)
	}
}

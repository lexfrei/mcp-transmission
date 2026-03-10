package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/go-transmission/api/transmission"
	"golang.org/x/sync/errgroup"

	"github.com/lexfrei/mcp-transmission/internal/config"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName        = "mcp-transmission"
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// version and revision are set via ldflags at build time.
var (
	version  = "dev"
	revision = "unknown"
)

func main() {
	err := run()
	if err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, cfgErr := config.Load()
	if cfgErr != nil {
		return errors.Wrap(cfgErr, "invalid configuration")
	}

	opts := []transmission.Option{}
	if cfg.HasAuth() {
		opts = append(opts, transmission.WithAuth(cfg.Username, cfg.Password))
	}

	client, err := transmission.New(cfg.TransmissionURL, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to create transmission client")
	}

	defer client.Close()

	serverOpts := newServerOptions()
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    serverName,
			Version: version + "+" + revision,
		},
		serverOpts,
	)

	registerTools(server, client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}

		signal.Stop(sigChan)
	}()

	group, groupCtx := errgroup.WithContext(ctx)
	httpEnabled := cfg.HTTPEnabled()

	group.Go(func() error {
		runErr := server.Run(groupCtx, &mcp.StdioTransport{})
		if runErr != nil && groupCtx.Err() == nil {
			return errors.Wrap(runErr, "stdio server failed")
		}

		// Only cancel when HTTP is not enabled; otherwise let the
		// HTTP transport keep running after stdin closes (e.g. in
		// container deployments without an interactive terminal).
		if !httpEnabled {
			cancel()
		}

		return nil
	})

	if httpEnabled {
		group.Go(func() error {
			return runHTTPServer(groupCtx, server, cfg.HTTPAddr())
		})
	}

	return group.Wait() //nolint:wrapcheck // errors are already wrapped inside group goroutines.
}

func newServerOptions() *mcp.ServerOptions {
	return &mcp.ServerOptions{
		Instructions: "MCP server for managing Transmission BitTorrent client. " +
			"Provides tools to list, add, remove, start, stop torrents, " +
			"view detailed torrent info, manage queue and bandwidth groups, " +
			"check session stats and configuration, test port accessibility, " +
			"and check free disk space. " +
			"Requires TRANSMISSION_URL environment variable " +
			"(defaults to http://localhost:9091/transmission/rpc). " +
			"Supports basic auth via TRANSMISSION_USERNAME/TRANSMISSION_PASSWORD.",
		Logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

func registerTools(server *mcp.Server, client transmission.Client) {
	mcp.AddTool(server, tools.TorrentListTool(), tools.NewTorrentListHandler(client))
	mcp.AddTool(server, tools.TorrentAddTool(), tools.NewTorrentAddHandler(client))
	mcp.AddTool(server, tools.TorrentRemoveTool(), tools.NewTorrentRemoveHandler(client))
	mcp.AddTool(server, tools.TorrentStartTool(), tools.NewTorrentStartHandler(client))
	mcp.AddTool(server, tools.TorrentStopTool(), tools.NewTorrentStopHandler(client))
	mcp.AddTool(server, tools.TorrentVerifyTool(), tools.NewTorrentVerifyHandler(client))
	mcp.AddTool(server, tools.TorrentReannounceTool(), tools.NewTorrentReannounceHandler(client))
	mcp.AddTool(server, tools.TorrentDetailsTool(), tools.NewTorrentDetailsHandler(client))
	mcp.AddTool(server, tools.TorrentSetTool(), tools.NewTorrentSetHandler(client))
	mcp.AddTool(server, tools.TorrentMoveTool(), tools.NewTorrentMoveHandler(client))
	mcp.AddTool(server, tools.SessionStatsTool(), tools.NewSessionStatsHandler(client))
	mcp.AddTool(server, tools.SessionGetTool(), tools.NewSessionGetHandler(client))
	mcp.AddTool(server, tools.SessionSetTool(), tools.NewSessionSetHandler(client))
	mcp.AddTool(server, tools.FreeSpaceTool(), tools.NewFreeSpaceHandler(client))
	mcp.AddTool(server, tools.PortTestTool(), tools.NewPortTestHandler(client))
	mcp.AddTool(server, tools.BlocklistUpdateTool(), tools.NewBlocklistUpdateHandler(client))
	mcp.AddTool(server, tools.QueueMoveTool(), tools.NewQueueMoveHandler(client))
	mcp.AddTool(server, tools.BandwidthGroupGetTool(), tools.NewBandwidthGroupGetHandler(client))
}

// runHTTPServer starts an HTTP/SSE transport for the MCP server.
// Sharing a single *mcp.Server across transports is safe: the SDK
// protects internal state (sessions, tools, subscriptions) with a sync.Mutex.
func runHTTPServer(ctx context.Context, server *mcp.Server, addr string) error {
	handler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server {
			return server
		},
		nil,
	)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	//nolint:gosec // G118: ctx is already cancelled when goroutine runs, must use fresh context for graceful shutdown.
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		shutdownErr := httpServer.Shutdown(shutdownCtx) //nolint:contextcheck // ctx is cancelled, need fresh context for graceful shutdown.
		if shutdownErr != nil {
			log.Printf("HTTP server shutdown error: %v", shutdownErr)
		}
	}()

	log.Printf("HTTP server listening on %s", addr)

	listenErr := httpServer.ListenAndServe()
	if errors.Is(listenErr, http.ErrServerClosed) {
		return nil
	}

	return errors.Wrap(listenErr, "HTTP listen failed")
}

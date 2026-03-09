package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSessionStatsTool_Definition(t *testing.T) {
	tool := tools.SessionStatsTool()

	if tool.Name != "transmission_session_stats" {
		t.Errorf("expected name transmission_session_stats, got %s", tool.Name)
	}
}

func TestSessionStatsHandler_Success(t *testing.T) {
	client := newMockClient()
	client.sessionStats = &transmission.SessionStats{
		ActiveTorrentCount: 3,
		PausedTorrentCount: 2,
		TorrentCount:       5,
		DownloadSpeed:      1024,
		UploadSpeed:        512,
	}

	handler := tools.NewSessionStatsHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionStatsParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.TotalTorrents != 5 {
		t.Errorf("expected 5 total torrents, got %d", output.TotalTorrents)
	}

	if output.ActiveTorrents != 3 {
		t.Errorf("expected 3 active torrents, got %d", output.ActiveTorrents)
	}
}

func TestSessionStatsHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewSessionStatsHandler(client)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionStatsParams{})
	if err == nil && (result == nil || !result.IsError) {
		t.Error("expected error")
	}
}

func TestSessionStatsHandler_NilResult(t *testing.T) {
	client := newMockClient()
	client.sessionStats = nil

	handler := tools.NewSessionStatsHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionStatsParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Output == "" {
		t.Error("expected non-empty output for nil stats")
	}
}

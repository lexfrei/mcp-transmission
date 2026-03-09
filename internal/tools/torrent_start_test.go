package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentStartTool_Definition(t *testing.T) {
	tool := tools.TorrentStartTool()

	if tool.Name != "transmission_torrent_start" {
		t.Errorf("expected name transmission_torrent_start, got %s", tool.Name)
	}
}

func TestTorrentStartHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{IDs: []int64{1}}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentStartHandler_Now(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{IDs: []int64{1}, Now: true}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentStartHandler_All(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{}

	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if output.Message != "Started all torrents" {
		t.Errorf("expected 'Started all torrents', got %s", output.Message)
	}
}

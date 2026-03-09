package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentReannounceTool_Definition(t *testing.T) {
	tool := tools.TorrentReannounceTool()

	if tool.Name != "transmission_torrent_reannounce" {
		t.Errorf("expected name transmission_torrent_reannounce, got %s", tool.Name)
	}
}

func TestTorrentReannounceHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentReannounceHandler(client)

	params := tools.TorrentReannounceParams{IDs: []int64{1}}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentReannounceHandler_EmptyIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentReannounceHandler(client)

	params := tools.TorrentReannounceParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty IDs, got: %v", err)
	}
}

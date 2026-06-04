package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentStopTool_Definition(t *testing.T) {
	tool := tools.TorrentStopTool()

	if tool.Name != "transmission_torrent_stop" {
		t.Errorf("expected name transmission_torrent_stop, got %s", tool.Name)
	}
}

func TestTorrentStopHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStopHandler(client)

	params := tools.TorrentStopParams{IDs: []int64{1}}

	result, out, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if out.Message != "Stopped 1 torrent(s)" {
		t.Errorf("expected message %q, got %q", "Stopped 1 torrent(s)", out.Message)
	}

	if client.LastMethod() != "TorrentStop" {
		t.Errorf("expected method TorrentStop, got %s", client.LastMethod())
	}
}

func TestTorrentStopHandler_EmptyIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStopHandler(client)

	params := tools.TorrentStopParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty IDs, got: %v", err)
	}
}

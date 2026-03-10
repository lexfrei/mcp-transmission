package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentVerifyTool_Definition(t *testing.T) {
	tool := tools.TorrentVerifyTool()

	if tool.Name != "transmission_torrent_verify" {
		t.Errorf("expected name transmission_torrent_verify, got %s", tool.Name)
	}
}

func TestTorrentVerifyHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentVerifyHandler(client)

	params := tools.TorrentVerifyParams{IDs: []int64{1}}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentVerifyHandler_EmptyIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentVerifyHandler(client)

	params := tools.TorrentVerifyParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty IDs, got: %v", err)
	}
}

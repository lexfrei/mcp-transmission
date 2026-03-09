package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentRemoveTool_Definition(t *testing.T) {
	tool := tools.TorrentRemoveTool()

	if tool.Name != "transmission_torrent_remove" {
		t.Errorf("expected name transmission_torrent_remove, got %s", tool.Name)
	}
}

func TestTorrentRemoveHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	params := tools.TorrentRemoveParams{IDs: []int64{1, 2}}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestTorrentRemoveHandler_MissingIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	params := tools.TorrentRemoveParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentRemoveHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock
	handler := tools.NewTorrentRemoveHandler(client)

	params := tools.TorrentRemoveParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

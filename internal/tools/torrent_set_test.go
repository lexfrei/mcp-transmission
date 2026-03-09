package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentSetTool_Definition(t *testing.T) {
	tool := tools.TorrentSetTool()

	if tool.Name != "transmission_torrent_set" {
		t.Errorf("expected name transmission_torrent_set, got %s", tool.Name)
	}
}

func TestTorrentSetHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentSetHandler(client)

	limit := int64(500)
	params := tools.TorrentSetParams{
		IDs:           []int64{1},
		DownloadLimit: &limit,
	}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentSetHandler_TransmissionError(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewTorrentSetHandler(client)

	limit := int64(500)
	params := tools.TorrentSetParams{
		IDs:           []int64{1},
		DownloadLimit: &limit,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

func TestTorrentSetHandler_EmptyChanges(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentSetHandler(client)

	params := tools.TorrentSetParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty changes, got: %v", err)
	}
}

func TestTorrentSetHandler_MissingIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentSetHandler(client)

	params := tools.TorrentSetParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

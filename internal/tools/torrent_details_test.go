package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentDetailsTool_Definition(t *testing.T) {
	tool := tools.TorrentDetailsTool()

	if tool.Name != "transmission_torrent_details" {
		t.Errorf("expected name transmission_torrent_details, got %s", tool.Name)
	}
}

func TestTorrentDetailsHandler_Success(t *testing.T) {
	client := newMockClient()
	name := "test-torrent"
	id := int64(1)
	status := transmission.TorrentStatusSeed

	client.torrentGetResult = &transmission.TorrentGetResult{
		Torrents: []transmission.Torrent{
			{ID: &id, Name: &name, Status: &status},
		},
	}

	handler := tools.NewTorrentDetailsHandler(client)

	params := tools.TorrentDetailsParams{ID: 1}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Output == "" {
		t.Error("expected non-empty output")
	}
}

func TestTorrentDetailsHandler_ZeroID(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentDetailsHandler(client)

	params := tools.TorrentDetailsParams{ID: 0}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentDetailsHandler_NotFound(t *testing.T) {
	client := newMockClient()
	client.torrentGetResult = &transmission.TorrentGetResult{
		Torrents: []transmission.Torrent{},
	}

	handler := tools.NewTorrentDetailsHandler(client)

	params := tools.TorrentDetailsParams{ID: 999}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentListTool_Definition(t *testing.T) {
	tool := tools.TorrentListTool()

	if tool.Name != "transmission_torrent_list" {
		t.Errorf("expected name transmission_torrent_list, got %s", tool.Name)
	}

	if tool.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestTorrentListHandler_Success(t *testing.T) {
	client := newMockClient()
	status := transmission.TorrentStatusSeed
	name := "test-torrent"
	id := int64(1)
	pct := 1.0

	client.torrentGetResult = &transmission.TorrentGetResult{
		Torrents: []transmission.Torrent{
			{
				ID:          &id,
				Name:        &name,
				Status:      &status,
				PercentDone: &pct,
			},
		},
	}

	handler := tools.NewTorrentListHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.TorrentListParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success, got error")
	}

	if output.Count != 1 {
		t.Errorf("expected count 1, got %d", output.Count)
	}
}

func TestTorrentListHandler_NilResult(t *testing.T) {
	client := newMockClient()
	client.torrentGetResult = nil

	handler := tools.NewTorrentListHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.TorrentListParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Count != 0 {
		t.Errorf("expected count 0, got %d", output.Count)
	}
}

func TestTorrentListHandler_Empty(t *testing.T) {
	client := newMockClient()
	client.torrentGetResult = &transmission.TorrentGetResult{
		Torrents: []transmission.Torrent{},
	}

	handler := tools.NewTorrentListHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.TorrentListParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success, got error")
	}

	if output.Count != 0 {
		t.Errorf("expected count 0, got %d", output.Count)
	}
}

func TestTorrentListHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewTorrentListHandler(client)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.TorrentListParams{})
	if err == nil && (result == nil || !result.IsError) {
		t.Error("expected error")
	}
}

package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentAddTool_Definition(t *testing.T) {
	tool := tools.TorrentAddTool()

	if tool.Name != "transmission_torrent_add" {
		t.Errorf("expected name transmission_torrent_add, got %s", tool.Name)
	}
}

func TestTorrentAddHandler_Success(t *testing.T) {
	client := newMockClient()
	client.torrentAddResult = &transmission.TorrentAddResult{
		TorrentAdded: &transmission.TorrentAddedInfo{
			ID:         1,
			Name:       "test-torrent",
			HashString: "abc123",
		},
	}

	handler := tools.NewTorrentAddHandler(client)

	params := tools.TorrentAddParams{
		Filename: "magnet:?xt=urn:btih:abc123",
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.ID != 1 {
		t.Errorf("expected ID 1, got %d", output.ID)
	}

	if output.Duplicate {
		t.Error("expected not duplicate")
	}
}

func TestTorrentAddHandler_Duplicate(t *testing.T) {
	client := newMockClient()
	client.torrentAddResult = &transmission.TorrentAddResult{
		TorrentDuplicate: &transmission.TorrentAddedInfo{
			ID:         1,
			Name:       "test-torrent",
			HashString: "abc123",
		},
	}

	handler := tools.NewTorrentAddHandler(client)

	params := tools.TorrentAddParams{
		Filename: "magnet:?xt=urn:btih:abc123",
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if !output.Duplicate {
		t.Error("expected duplicate")
	}
}

func TestTorrentAddHandler_MissingParams(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentAddHandler(client)

	params := tools.TorrentAddParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err == nil {
		t.Fatal("expected error for missing params")
	}

	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

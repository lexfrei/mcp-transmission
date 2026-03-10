package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentMoveTool_Definition(t *testing.T) {
	tool := tools.TorrentMoveTool()

	if tool.Name != "transmission_torrent_move" {
		t.Errorf("expected name transmission_torrent_move, got %s", tool.Name)
	}
}

func TestTorrentMoveHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentMoveHandler(client)

	params := tools.TorrentMoveParams{
		IDs:      []int64{1},
		Location: "/new/path",
		Move:     true,
	}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestTorrentMoveHandler_TransmissionError(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewTorrentMoveHandler(client)

	params := tools.TorrentMoveParams{
		IDs:      []int64{1},
		Location: "/new/path",
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

func TestTorrentMoveHandler_MissingIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentMoveHandler(client)

	params := tools.TorrentMoveParams{Location: "/path"}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentMoveHandler_MissingLocation(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentMoveHandler(client)

	params := tools.TorrentMoveParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentMoveHandler_RelativePath(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentMoveHandler(client)

	params := tools.TorrentMoveParams{
		IDs:      []int64{1},
		Location: "relative/path",
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for relative path, got: %v", err)
	}
}

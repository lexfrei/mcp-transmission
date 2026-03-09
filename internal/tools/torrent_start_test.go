package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
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

	if client.LastMethod() != "TorrentStart" {
		t.Errorf("expected TorrentStart, got %s", client.LastMethod())
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

	if client.LastMethod() != "TorrentStartNow" {
		t.Errorf("expected TorrentStartNow, got %s", client.LastMethod())
	}
}

func TestTorrentStartHandler_EmptyIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty IDs, got: %v", err)
	}
}

func TestTorrentStartHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

func TestTorrentStartHandler_NowError(t *testing.T) {
	client := newMockClient()
	client.err = errMock
	handler := tools.NewTorrentStartHandler(client)

	params := tools.TorrentStartParams{IDs: []int64{1}, Now: true}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

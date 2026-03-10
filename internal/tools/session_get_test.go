package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSessionGetTool_Definition(t *testing.T) {
	tool := tools.SessionGetTool()

	if tool.Name != "transmission_session_get" {
		t.Errorf("expected name transmission_session_get, got %s", tool.Name)
	}
}

func TestSessionGetHandler_Success(t *testing.T) {
	client := newMockClient()
	version := "4.0.0"
	downloadDir := "/downloads"

	client.sessionResult = &transmission.Session{
		Version:     &version,
		DownloadDir: &downloadDir,
	}

	handler := tools.NewSessionGetHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionGetParams{})
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

func TestSessionGetHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewSessionGetHandler(client)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionGetParams{})
	if err == nil && (result == nil || !result.IsError) {
		t.Error("expected error")
	}
}

func TestSessionGetHandler_NilResult(t *testing.T) {
	client := newMockClient()
	client.sessionResult = nil

	handler := tools.NewSessionGetHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.SessionGetParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Output == "" {
		t.Error("expected non-empty output for nil session")
	}
}

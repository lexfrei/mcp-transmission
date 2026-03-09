package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestFreeSpaceTool_Definition(t *testing.T) {
	tool := tools.FreeSpaceTool()

	if tool.Name != "transmission_free_space" {
		t.Errorf("expected name transmission_free_space, got %s", tool.Name)
	}
}

func TestFreeSpaceHandler_Success(t *testing.T) {
	client := newMockClient()
	client.freeSpaceResult = &transmission.FreeSpace{
		Path:      "/downloads",
		SizeBytes: 107374182400,
	}

	handler := tools.NewFreeSpaceHandler(client)

	params := tools.FreeSpaceParams{Path: "/downloads"}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Path != "/downloads" {
		t.Errorf("expected path /downloads, got %s", output.Path)
	}
}

func TestFreeSpaceHandler_MissingPath(t *testing.T) {
	client := newMockClient()
	handler := tools.NewFreeSpaceHandler(client)

	params := tools.FreeSpaceParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestFreeSpaceHandler_RelativePath(t *testing.T) {
	client := newMockClient()
	handler := tools.NewFreeSpaceHandler(client)

	params := tools.FreeSpaceParams{Path: "relative/dir"}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for relative path, got: %v", err)
	}
}

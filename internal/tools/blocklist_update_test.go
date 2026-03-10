package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestBlocklistUpdateTool_Definition(t *testing.T) {
	tool := tools.BlocklistUpdateTool()

	if tool.Name != "transmission_blocklist_update" {
		t.Errorf("expected name transmission_blocklist_update, got %s", tool.Name)
	}
}

func TestBlocklistUpdateHandler_Success(t *testing.T) {
	client := newMockClient()
	client.blocklistCount = 42000

	handler := tools.NewBlocklistUpdateHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.BlocklistUpdateParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.BlockedRanges != 42000 {
		t.Errorf("expected 42000 blocked ranges, got %d", output.BlockedRanges)
	}
}

package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestBandwidthGroupGetTool_Definition(t *testing.T) {
	tool := tools.BandwidthGroupGetTool()

	if tool.Name != "transmission_bandwidth_group_get" {
		t.Errorf("expected name transmission_bandwidth_group_get, got %s", tool.Name)
	}
}

func TestBandwidthGroupGetHandler_Success(t *testing.T) {
	client := newMockClient()
	client.bandwidthGroups = []transmission.BandwidthGroup{
		{Name: "default"},
	}

	handler := tools.NewBandwidthGroupGetHandler(client)

	params := tools.BandwidthGroupGetParams{}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Count != 1 {
		t.Errorf("expected count 1, got %d", output.Count)
	}
}

func TestBandwidthGroupGetHandler_Empty(t *testing.T) {
	client := newMockClient()
	client.bandwidthGroups = []transmission.BandwidthGroup{}

	handler := tools.NewBandwidthGroupGetHandler(client)

	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.BandwidthGroupGetParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if output.Count != 0 {
		t.Errorf("expected count 0, got %d", output.Count)
	}
}

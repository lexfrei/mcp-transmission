package tools_test

import (
	"context"
	"testing"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestPortTestTool_Definition(t *testing.T) {
	tool := tools.PortTestTool()

	if tool.Name != "transmission_port_test" {
		t.Errorf("expected name transmission_port_test, got %s", tool.Name)
	}
}

func TestPortTestHandler_Open(t *testing.T) {
	client := newMockClient()
	client.portTestResult = true

	handler := tools.NewPortTestHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.PortTestParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if !output.PortOpen {
		t.Error("expected port open")
	}
}

func TestPortTestHandler_Closed(t *testing.T) {
	client := newMockClient()
	client.portTestResult = false

	handler := tools.NewPortTestHandler(client)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, tools.PortTestParams{})
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.PortOpen {
		t.Error("expected port closed")
	}
}

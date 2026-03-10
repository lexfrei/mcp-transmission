package tools

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// PortTestParams defines the parameters for the transmission_port_test tool.
type PortTestParams struct{}

// PortTestResult is the output of the transmission_port_test tool.
type PortTestResult struct {
	PortOpen bool   `json:"portOpen"`
	Message  string `json:"message"`
}

// NewPortTestHandler creates a handler for the transmission_port_test tool.
func NewPortTestHandler(client transmission.Client) mcp.ToolHandlerFor[PortTestParams, PortTestResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ PortTestParams,
	) (*mcp.CallToolResult, PortTestResult, error) {
		portOpen, err := client.PortTest(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, PortTestResult{},
				transmissionErr("port test failed", err)
		}

		msg := "Peer port is open and reachable"
		if !portOpen {
			msg = "Peer port is NOT reachable from the outside"
		}

		return nil, PortTestResult{PortOpen: portOpen, Message: msg}, nil
	}
}

// PortTestTool returns the MCP tool definition for transmission_port_test.
func PortTestTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_port_test",
		Description: "Test if the Transmission peer port is accessible from the internet",
	}
}

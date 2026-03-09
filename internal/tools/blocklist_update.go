package tools

import (
	"context"
	"fmt"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BlocklistUpdateParams defines the parameters for the transmission_blocklist_update tool.
type BlocklistUpdateParams struct{}

// BlocklistUpdateResult is the output of the transmission_blocklist_update tool.
type BlocklistUpdateResult struct {
	BlockedRanges int    `json:"blockedRanges"`
	Message       string `json:"message"`
}

// NewBlocklistUpdateHandler creates a handler for the transmission_blocklist_update tool.
func NewBlocklistUpdateHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[BlocklistUpdateParams, BlocklistUpdateResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ BlocklistUpdateParams,
	) (*mcp.CallToolResult, BlocklistUpdateResult, error) {
		count, err := client.BlocklistUpdate(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, BlocklistUpdateResult{},
				transmissionErr("failed to update blocklist", err)
		}

		return nil, BlocklistUpdateResult{
			BlockedRanges: count,
			Message:       fmt.Sprintf("Blocklist updated: %d IP ranges blocked", count),
		}, nil
	}
}

// BlocklistUpdateTool returns the MCP tool definition for transmission_blocklist_update.
func BlocklistUpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_blocklist_update",
		Description: "Update the IP blocklist from the configured URL",
	}
}

package tools

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentVerifyParams defines the parameters for the transmission_torrent_verify tool.
type TorrentVerifyParams struct {
	IDs []int64 `json:"ids,omitempty" jsonschema:"Torrent IDs to verify (empty = all)"`
}

// TorrentVerifyResult is the output of the transmission_torrent_verify tool.
type TorrentVerifyResult struct {
	Message string `json:"message"`
}

// NewTorrentVerifyHandler creates a handler for the transmission_torrent_verify tool.
func NewTorrentVerifyHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentVerifyParams, TorrentVerifyResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentVerifyParams,
	) (*mcp.CallToolResult, TorrentVerifyResult, error) {
		verifyErr := client.TorrentVerify(ctx, params.IDs)
		if verifyErr != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentVerifyResult{},
				transmissionErr("failed to verify torrents", verifyErr)
		}

		return nil, TorrentVerifyResult{
			Message: formatActionMessage("Queued verification for", params.IDs),
		}, nil
	}
}

// TorrentVerifyTool returns the MCP tool definition for transmission_torrent_verify.
func TorrentVerifyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_verify",
		Description: "Verify local data integrity for one or more torrents",
	}
}

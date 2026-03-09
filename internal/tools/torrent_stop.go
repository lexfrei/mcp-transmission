package tools

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentStopParams defines the parameters for the transmission_torrent_stop tool.
type TorrentStopParams struct {
	IDs []int64 `json:"ids" jsonschema:"Torrent IDs to stop"`
}

// TorrentStopResult is the output of the transmission_torrent_stop tool.
type TorrentStopResult struct {
	Message string `json:"message"`
}

// NewTorrentStopHandler creates a handler for the transmission_torrent_stop tool.
func NewTorrentStopHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentStopParams, TorrentStopResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentStopParams,
	) (*mcp.CallToolResult, TorrentStopResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentStopResult{},
				validationErr(ErrIDsRequired)
		}

		err := client.TorrentStop(ctx, params.IDs)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentStopResult{},
				transmissionErr("failed to stop torrents", err)
		}

		return nil, TorrentStopResult{
			Message: formatActionMessage("Stopped", params.IDs),
		}, nil
	}
}

// TorrentStopTool returns the MCP tool definition for transmission_torrent_stop.
func TorrentStopTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_stop",
		Description: "Stop (pause) one or more torrents",
	}
}

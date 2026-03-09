package tools

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentReannounceParams defines the parameters for the transmission_torrent_reannounce tool.
type TorrentReannounceParams struct {
	IDs []int64 `json:"ids" jsonschema:"Torrent IDs to reannounce"`
}

// TorrentReannounceResult is the output of the transmission_torrent_reannounce tool.
type TorrentReannounceResult struct {
	Message string `json:"message"`
}

// NewTorrentReannounceHandler creates a handler for the transmission_torrent_reannounce tool.
func NewTorrentReannounceHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentReannounceParams, TorrentReannounceResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentReannounceParams,
	) (*mcp.CallToolResult, TorrentReannounceResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentReannounceResult{},
				validationErr(ErrIDsRequired)
		}

		announceErr := client.TorrentReannounce(ctx, params.IDs)
		if announceErr != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentReannounceResult{},
				transmissionErr("failed to reannounce torrents", announceErr)
		}

		return nil, TorrentReannounceResult{
			Message: formatActionMessage("Reannounced", params.IDs),
		}, nil
	}
}

// TorrentReannounceTool returns the MCP tool definition for transmission_torrent_reannounce.
func TorrentReannounceTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_reannounce",
		Description: "Force immediate tracker announce for one or more torrents",
	}
}

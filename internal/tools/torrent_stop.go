package tools

import (
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentStopParams defines the parameters for the transmission_torrent_stop tool.
type TorrentStopParams struct {
	IDs []int64 `json:"ids" jsonschema:"Torrent IDs to stop"`
}

func (p TorrentStopParams) torrentIDs() []int64 { return p.IDs }

// NewTorrentStopHandler creates a handler for the transmission_torrent_stop tool.
func NewTorrentStopHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentStopParams, idsResult] {
	return newIDsActionHandler[TorrentStopParams](
		client.TorrentStop,
		"failed to stop torrents",
		"Stopped",
	)
}

// TorrentStopTool returns the MCP tool definition for transmission_torrent_stop.
func TorrentStopTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_stop",
		Description: "Stop (pause) one or more torrents",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Stop Torrents",
			IdempotentHint:  true,
			DestructiveHint: ptrBool(false),
		},
	}
}

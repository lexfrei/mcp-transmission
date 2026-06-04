package tools

import (
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentReannounceParams defines the parameters for the transmission_torrent_reannounce tool.
type TorrentReannounceParams struct {
	IDs []int64 `json:"ids" jsonschema:"Torrent IDs to reannounce"`
}

func (p TorrentReannounceParams) torrentIDs() []int64 { return p.IDs }

// NewTorrentReannounceHandler creates a handler for the transmission_torrent_reannounce tool.
func NewTorrentReannounceHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentReannounceParams, idsResult] {
	return newIDsActionHandler[TorrentReannounceParams](
		client.TorrentReannounce,
		"failed to reannounce torrents",
		"Reannounced",
	)
}

// TorrentReannounceTool returns the MCP tool definition for transmission_torrent_reannounce.
func TorrentReannounceTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_reannounce",
		Description: "Force immediate tracker announce for one or more torrents",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Reannounce Torrents",
			IdempotentHint:  true,
			DestructiveHint: ptrBool(false),
		},
	}
}

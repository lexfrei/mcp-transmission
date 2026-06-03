package tools

import (
	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentVerifyParams defines the parameters for the transmission_torrent_verify tool.
type TorrentVerifyParams struct {
	IDs []int64 `json:"ids" jsonschema:"Torrent IDs to verify"`
}

func (p TorrentVerifyParams) torrentIDs() []int64 { return p.IDs }

// NewTorrentVerifyHandler creates a handler for the transmission_torrent_verify tool.
func NewTorrentVerifyHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentVerifyParams, idsResult] {
	return newIDsActionHandler[TorrentVerifyParams](
		client.TorrentVerify,
		"failed to verify torrents",
		"Queued verification for",
	)
}

// TorrentVerifyTool returns the MCP tool definition for transmission_torrent_verify.
func TorrentVerifyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_verify",
		Description: "Verify local data integrity for one or more torrents",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Verify Torrents",
			IdempotentHint:  true,
			DestructiveHint: ptrBool(false),
		},
	}
}

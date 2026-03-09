package tools

import (
	"context"
	"fmt"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TorrentStartParams defines the parameters for the transmission_torrent_start tool.
type TorrentStartParams struct {
	IDs []int64 `json:"ids"           jsonschema:"Torrent IDs to start"`
	Now bool    `json:"now,omitempty" jsonschema:"Start immediately, bypassing the queue"`
}

// TorrentStartResult is the output of the transmission_torrent_start tool.
type TorrentStartResult struct {
	Message string `json:"message"`
}

// NewTorrentStartHandler creates a handler for the transmission_torrent_start tool.
func NewTorrentStartHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentStartParams, TorrentStartResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentStartParams,
	) (*mcp.CallToolResult, TorrentStartResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentStartResult{},
				validationErr(ErrIDsRequired)
		}

		var err error

		if params.Now {
			err = client.TorrentStartNow(ctx, params.IDs)
		} else {
			err = client.TorrentStart(ctx, params.IDs)
		}

		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentStartResult{},
				transmissionErr("failed to start torrents", err)
		}

		msg := formatActionMessage("Started", params.IDs)

		return nil, TorrentStartResult{Message: msg}, nil
	}
}

// TorrentStartTool returns the MCP tool definition for transmission_torrent_start.
func TorrentStartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_start",
		Description: "Start one or more torrents. Use 'now' to bypass queue",
	}
}

func formatActionMessage(action string, ids []int64) string {
	return fmt.Sprintf("%s %d torrent(s)", action, len(ids))
}

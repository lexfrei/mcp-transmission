package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrLocationRequired is returned when the location parameter is missing.
var ErrLocationRequired = errors.New("location is required")

// TorrentMoveParams defines the parameters for the transmission_torrent_move tool.
type TorrentMoveParams struct {
	IDs      []int64 `json:"ids"      jsonschema:"Torrent IDs to move"`
	Location string  `json:"location" jsonschema:"New download directory path"`
	Move     bool    `json:"move"     jsonschema:"Move existing files to new location (if false, only update path)"`
}

// TorrentMoveResult is the output of the transmission_torrent_move tool.
type TorrentMoveResult struct {
	Message string `json:"message"`
}

// NewTorrentMoveHandler creates a handler for the transmission_torrent_move tool.
func NewTorrentMoveHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentMoveParams, TorrentMoveResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentMoveParams,
	) (*mcp.CallToolResult, TorrentMoveResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentMoveResult{},
				validationErr(ErrIDsRequired)
		}

		if params.Location == "" {
			return &mcp.CallToolResult{IsError: true}, TorrentMoveResult{},
				validationErr(ErrLocationRequired)
		}

		if !filepath.IsAbs(params.Location) {
			return &mcp.CallToolResult{IsError: true}, TorrentMoveResult{},
				validationErr(ErrAbsolutePathRequired)
		}

		err := client.TorrentSetLocation(ctx, params.IDs, params.Location, params.Move)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentMoveResult{},
				transmissionErr("failed to move torrents", err)
		}

		action := "Updated path for"
		if params.Move {
			action = "Moving files for"
		}

		msg := fmt.Sprintf("%s %d torrent(s) to %s", action, len(params.IDs), params.Location)

		return nil, TorrentMoveResult{Message: msg}, nil
	}
}

// TorrentMoveTool returns the MCP tool definition for transmission_torrent_move.
func TorrentMoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_move",
		Description: "Move torrent data to a new location on disk",
	}
}

package tools

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrIDsRequired is returned when the ids parameter is missing.
var ErrIDsRequired = errors.New("at least one torrent ID is required")

// ErrDeleteConfirmRequired is returned when deleteLocalData is true but confirmDelete is not.
var ErrDeleteConfirmRequired = errors.New("confirmDelete must be true when deleteLocalData is true")

// TorrentRemoveParams defines the parameters for the transmission_torrent_remove tool.
type TorrentRemoveParams struct {
	IDs             []int64 `json:"ids"                       jsonschema:"Torrent IDs to remove"`
	DeleteLocalData bool    `json:"deleteLocalData,omitempty" jsonschema:"Also delete downloaded files (DESTRUCTIVE)"`
	ConfirmDelete   bool    `json:"confirmDelete,omitempty"   jsonschema:"Must be true when deleteLocalData is true"`
}

// TorrentRemoveResult is the output of the transmission_torrent_remove tool.
type TorrentRemoveResult struct {
	Message string `json:"message"`
}

// NewTorrentRemoveHandler creates a handler for the transmission_torrent_remove tool.
func NewTorrentRemoveHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentRemoveParams,
	) (*mcp.CallToolResult, TorrentRemoveResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
				validationErr(ErrIDsRequired)
		}

		if params.DeleteLocalData && !params.ConfirmDelete {
			return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
				validationErr(ErrDeleteConfirmRequired)
		}

		err := client.TorrentRemove(ctx, params.IDs, params.DeleteLocalData)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
				transmissionErr("failed to remove torrents", err)
		}

		msg := fmt.Sprintf("Removed %d torrent(s)", len(params.IDs))
		if params.DeleteLocalData {
			msg += " and their local data"
		}

		return nil, TorrentRemoveResult{Message: msg}, nil
	}
}

// TorrentRemoveTool returns the MCP tool definition for transmission_torrent_remove.
func TorrentRemoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_remove",
		Description: "Remove one or more torrents. Set deleteLocalData=true AND confirmDelete=true to also delete files from disk (DESTRUCTIVE)",
	}
}

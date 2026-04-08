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

// ErrElicitationFailed is returned when the elicitation request fails.
var ErrElicitationFailed = errors.New("failed to request deletion confirmation")

// ErrNoElicitation is returned when no elicitation support is available.
var ErrNoElicitation = errors.New("no elicitation support available")

// ErrNilElicitResult is returned when the elicitor returns a nil result without an error.
var ErrNilElicitResult = errors.New("elicitation returned nil result")

// Elicitor abstracts the MCP elicitation capability for testability.
type Elicitor interface {
	Elicit(ctx context.Context, params *mcp.ElicitParams) (*mcp.ElicitResult, error)
}

// TorrentRemoveParams defines the parameters for the transmission_torrent_remove tool.
type TorrentRemoveParams struct {
	IDs             []int64 `json:"ids"                       jsonschema:"Torrent IDs to remove"`
	DeleteLocalData bool    `json:"deleteLocalData,omitempty" jsonschema:"Also delete downloaded files (DESTRUCTIVE, requires confirmation)"`
}

// TorrentRemoveResult is the output of the transmission_torrent_remove tool.
type TorrentRemoveResult struct {
	Message string `json:"message"`
}

// NewTorrentRemoveHandler creates a handler for the transmission_torrent_remove tool.
// It uses the MCP session from the request for elicitation when deleteLocalData is true.
func NewTorrentRemoveHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return newTorrentRemoveHandler(client, nil)
}

// NewTorrentRemoveHandlerWithElicitor creates a handler with an explicit Elicitor for testing.
func NewTorrentRemoveHandlerWithElicitor(client transmission.Client, elicitor Elicitor) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return newTorrentRemoveHandler(client, elicitor)
}

func newTorrentRemoveHandler(client transmission.Client, elicitor Elicitor) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		params TorrentRemoveParams,
	) (*mcp.CallToolResult, TorrentRemoveResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
				validationErr(ErrIDsRequired)
		}

		if params.DeleteLocalData {
			confirmed, confirmErr := confirmDeletion(ctx, req, elicitor, params.IDs)
			if confirmErr != nil {
				return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
					errors.Mark(errors.Wrap(confirmErr, "elicitation failed"), ErrElicitationFailed)
			}

			if !confirmed {
				return nil, TorrentRemoveResult{
					Message: "Deletion cancelled by user",
				}, nil
			}
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

func confirmDeletion(
	ctx context.Context,
	req *mcp.CallToolRequest,
	elicitor Elicitor,
	ids []int64,
) (bool, error) {
	resolvedElicitor := elicitor
	if resolvedElicitor == nil && req != nil && req.Session != nil {
		resolvedElicitor = req.Session
	}

	if resolvedElicitor == nil {
		return false, ErrNoElicitation
	}

	result, elicitErr := resolvedElicitor.Elicit(ctx, &mcp.ElicitParams{
		Message: fmt.Sprintf(
			"Are you sure you want to permanently delete local data for %d torrent(s)? This action cannot be undone.",
			len(ids),
		),
	})
	if elicitErr != nil {
		return false, errors.Wrap(elicitErr, "elicit call failed")
	}

	if result == nil {
		return false, ErrNilElicitResult
	}

	return result.Action == "accept", nil
}

// TorrentRemoveTool returns the MCP tool definition for transmission_torrent_remove.
func TorrentRemoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_remove",
		Description: "Remove one or more torrents. Set deleteLocalData=true to also delete files from disk (DESTRUCTIVE, requires confirmation via elicitation)",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Remove Torrents",
			DestructiveHint: ptrBool(true),
			// Idempotent: removing an already-removed torrent is a no-op in Transmission.
			// Even with deleteLocalData=true, the second call has no effect (files already gone).
			IdempotentHint: true,
		},
	}
}

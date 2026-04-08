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

// ErrDeleteConfirmRequired is returned when deleteLocalData is true but neither
// elicitation nor confirmDelete confirmed the operation.
var ErrDeleteConfirmRequired = errors.New(
	"deleteLocalData requires confirmation: " +
		"either use a client that supports elicitation, " +
		"or set confirmDelete=true",
)

var (
	errNoElicitation   = errors.New("no elicitation support available")
	errNilElicitResult = errors.New("elicitation returned nil result")
)

// Elicitor abstracts the MCP elicitation capability for testability.
type Elicitor interface {
	Elicit(ctx context.Context, params *mcp.ElicitParams) (*mcp.ElicitResult, error)
}

// TorrentRemoveParams defines the parameters for the transmission_torrent_remove tool.
type TorrentRemoveParams struct {
	IDs             []int64 `json:"ids"                       jsonschema:"Torrent IDs to remove"`
	DeleteLocalData bool    `json:"deleteLocalData,omitempty" jsonschema:"Also delete downloaded files (DESTRUCTIVE, requires confirmation)"`
	ConfirmDelete   bool    `json:"confirmDelete,omitempty"   jsonschema:"Fallback confirmation when client lacks elicitation support"`
}

// TorrentRemoveResult is the output of the transmission_torrent_remove tool.
type TorrentRemoveResult struct {
	Message string `json:"message"`
}

// NewTorrentRemoveHandler creates a handler for the transmission_torrent_remove tool.
// It uses the MCP session from the request for elicitation when deleteLocalData is true.
// Falls back to the confirmDelete parameter if elicitation is unavailable.
func NewTorrentRemoveHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return newTorrentRemoveHandler(client, nil)
}

// NewTorrentRemoveHandlerWithElicitor creates a handler with an explicit Elicitor for testing.
func NewTorrentRemoveHandlerWithElicitor(
	client transmission.Client,
	elicitor Elicitor,
) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
	return newTorrentRemoveHandler(client, elicitor)
}

func newTorrentRemoveHandler(
	client transmission.Client,
	elicitor Elicitor,
) mcp.ToolHandlerFor[TorrentRemoveParams, TorrentRemoveResult] {
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
			confirmed, confirmErr := confirmDeletion(ctx, req, elicitor, &params)
			if confirmErr != nil {
				return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
					confirmErr
			}

			if !confirmed {
				return nil, TorrentRemoveResult{
					Message: "Deletion cancelled by user",
				}, nil
			}
		}

		removeErr := client.TorrentRemove(ctx, params.IDs, params.DeleteLocalData)
		if removeErr != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentRemoveResult{},
				transmissionErr("failed to remove torrents", removeErr)
		}

		msg := fmt.Sprintf("Removed %d torrent(s)", len(params.IDs))
		if params.DeleteLocalData {
			msg += " and their local data"
		}

		return nil, TorrentRemoveResult{Message: msg}, nil
	}
}

// confirmDeletion tries elicitation first. If elicitation is unavailable (no session,
// no elicitor, or the client returns an error), it falls back to the confirmDelete parameter.
func confirmDeletion(
	ctx context.Context,
	req *mcp.CallToolRequest,
	elicitor Elicitor,
	params *TorrentRemoveParams,
) (bool, error) {
	confirmed, elicitErr := tryElicit(ctx, req, elicitor, params.IDs)
	if elicitErr == nil {
		return confirmed, nil
	}

	// Elicitation unavailable or failed — fall back to confirmDelete param.
	if params.ConfirmDelete {
		return true, nil
	}

	return false, validationErr(ErrDeleteConfirmRequired)
}

func tryElicit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	elicitor Elicitor,
	ids []int64,
) (bool, error) {
	resolved := elicitor
	if resolved == nil && req != nil && req.Session != nil {
		resolved = req.Session
	}

	if resolved == nil {
		return false, errNoElicitation
	}

	result, elicitErr := resolved.Elicit(ctx, &mcp.ElicitParams{
		Message: fmt.Sprintf(
			"Are you sure you want to permanently delete local data for %d torrent(s)? This action cannot be undone.",
			len(ids),
		),
	})
	if elicitErr != nil {
		return false, errors.Wrap(elicitErr, "elicit call failed")
	}

	if result == nil {
		return false, errNilElicitResult
	}

	return result.Action == "accept", nil
}

// TorrentRemoveTool returns the MCP tool definition for transmission_torrent_remove.
func TorrentRemoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "transmission_torrent_remove",
		Description: "Remove one or more torrents. Set deleteLocalData=true to also delete files from disk (DESTRUCTIVE). " +
			"Confirmation is requested via elicitation if the client supports it; " +
			"otherwise set confirmDelete=true as a fallback.",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Remove Torrents",
			DestructiveHint: ptrBool(true),
			// Idempotent: removing an already-removed torrent is a no-op in Transmission RPC.
			// Even with deleteLocalData=true, the second call has no effect (files already gone).
			// Note: elicitation-capable clients will see a confirmation prompt on each call,
			// but the underlying operation remains idempotent.
			IdempotentHint: true,
		},
	}
}

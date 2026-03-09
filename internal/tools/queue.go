package tools

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Valid queue actions.
const (
	QueueActionTop    = "top"
	QueueActionUp     = "up"
	QueueActionDown   = "down"
	QueueActionBottom = "bottom"
)

// ErrInvalidQueueAction is returned when the action parameter is not valid.
var ErrInvalidQueueAction = errors.New("action must be one of: top, up, down, bottom")

// QueueMoveParams defines the parameters for the transmission_queue_move tool.
type QueueMoveParams struct {
	IDs    []int64 `json:"ids"    jsonschema:"Torrent IDs to move in queue"`
	Action string  `json:"action" jsonschema:"Queue action: top, up, down, or bottom"`
}

// QueueMoveResult is the output of the transmission_queue_move tool.
type QueueMoveResult struct {
	Message string `json:"message"`
}

// NewQueueMoveHandler creates a handler for the transmission_queue_move tool.
func NewQueueMoveHandler(client transmission.Client) mcp.ToolHandlerFor[QueueMoveParams, QueueMoveResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params QueueMoveParams,
	) (*mcp.CallToolResult, QueueMoveResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, QueueMoveResult{},
				validationErr(ErrIDsRequired)
		}

		err := executeQueueAction(ctx, client, params.IDs, params.Action)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, QueueMoveResult{}, err
		}

		msg := fmt.Sprintf("Moved %d torrent(s) to queue %s", len(params.IDs), params.Action)

		return nil, QueueMoveResult{Message: msg}, nil
	}
}

// QueueMoveTool returns the MCP tool definition for transmission_queue_move.
func QueueMoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_queue_move",
		Description: "Move torrents in the download queue (top, up, down, bottom)",
	}
}

func executeQueueAction(
	ctx context.Context,
	client transmission.Client,
	ids []int64,
	action string,
) error {
	var err error

	switch action {
	case QueueActionTop:
		err = client.QueueMoveTop(ctx, ids)
	case QueueActionUp:
		err = client.QueueMoveUp(ctx, ids)
	case QueueActionDown:
		err = client.QueueMoveDown(ctx, ids)
	case QueueActionBottom:
		err = client.QueueMoveBottom(ctx, ids)
	default:
		return validationErr(ErrInvalidQueueAction)
	}

	if err != nil {
		return transmissionErr("queue move failed", err)
	}

	return nil
}

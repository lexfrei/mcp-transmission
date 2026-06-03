package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// idsResult is the common output shape for tools that act on a list of
// torrent IDs and report back a single human-readable message.
type idsResult struct {
	Message string `json:"message"`
}

// idsAction performs an action against the Transmission RPC API for a set of
// torrent IDs.
type idsAction func(ctx context.Context, ids []int64) error

// idsParams is implemented by parameter structs that carry a list of torrent
// IDs, allowing a single generic handler to validate and forward them.
type idsParams interface {
	torrentIDs() []int64
}

// newIDsActionHandler builds an MCP handler for tools whose only behaviour is
// to validate a non-empty list of torrent IDs, invoke a single action against
// the client, and return a formatted success message. The errMsg is used to
// wrap any transport failure, and successMsg is the verb reported on success.
func newIDsActionHandler[P idsParams](
	action idsAction,
	errMsg string,
	successMsg string,
) mcp.ToolHandlerFor[P, idsResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params P,
	) (*mcp.CallToolResult, idsResult, error) {
		ids := params.torrentIDs()
		if len(ids) == 0 {
			return &mcp.CallToolResult{IsError: true}, idsResult{},
				validationErr(ErrIDsRequired)
		}

		err := action(ctx, ids)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, idsResult{},
				transmissionErr(errMsg, err)
		}

		return nil, idsResult{Message: formatActionMessage(successMsg, ids)}, nil
	}
}

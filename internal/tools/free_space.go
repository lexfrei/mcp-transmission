package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrPathRequired is returned when the path parameter is missing.
var ErrPathRequired = errors.New("path is required")

// FreeSpaceParams defines the parameters for the transmission_free_space tool.
type FreeSpaceParams struct {
	Path string `json:"path" jsonschema:"Directory path to check free space for"`
}

// FreeSpaceResult is the output of the transmission_free_space tool.
type FreeSpaceResult struct {
	Path      string `json:"path"`
	FreeBytes int64  `json:"freeBytes"`
	FreeHuman string `json:"freeHuman"`
	Output    string `json:"output"`
}

// NewFreeSpaceHandler creates a handler for the transmission_free_space tool.
func NewFreeSpaceHandler(client transmission.Client) mcp.ToolHandlerFor[FreeSpaceParams, FreeSpaceResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params FreeSpaceParams,
	) (*mcp.CallToolResult, FreeSpaceResult, error) {
		if params.Path == "" {
			return &mcp.CallToolResult{IsError: true}, FreeSpaceResult{},
				validationErr(ErrPathRequired)
		}

		if !strings.HasPrefix(params.Path, "/") {
			return &mcp.CallToolResult{IsError: true}, FreeSpaceResult{},
				validationErr(ErrAbsolutePathRequired)
		}

		space, err := client.FreeSpace(ctx, params.Path)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, FreeSpaceResult{},
				transmissionErr("failed to check free space", err)
		}

		if space == nil {
			return nil, FreeSpaceResult{Output: "No free space data returned"}, nil
		}

		freeHuman := formatBytes(space.SizeBytes)

		result := FreeSpaceResult{
			Path:      space.Path,
			FreeBytes: space.SizeBytes,
			FreeHuman: freeHuman,
			Output:    fmt.Sprintf("Free space at %s: %s", space.Path, freeHuman),
		}

		return nil, result, nil
	}
}

// FreeSpaceTool returns the MCP tool definition for transmission_free_space.
func FreeSpaceTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_free_space",
		Description: "Check available disk space at a given path on the Transmission server",
	}
}

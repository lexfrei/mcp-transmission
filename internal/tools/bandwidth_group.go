package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BandwidthGroupGetParams defines the parameters for the transmission_bandwidth_group_get tool.
type BandwidthGroupGetParams struct {
	Names []string `json:"names,omitempty" jsonschema:"Group names to retrieve (empty = all)"`
}

// BandwidthGroupGetResult is the output of the transmission_bandwidth_group_get tool.
type BandwidthGroupGetResult struct {
	Count  int    `json:"count"`
	Output string `json:"output"`
}

// NewBandwidthGroupGetHandler creates a handler for the transmission_bandwidth_group_get tool.
func NewBandwidthGroupGetHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[BandwidthGroupGetParams, BandwidthGroupGetResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params BandwidthGroupGetParams,
	) (*mcp.CallToolResult, BandwidthGroupGetResult, error) {
		groups, err := client.GroupGet(ctx, params.Names)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, BandwidthGroupGetResult{},
				transmissionErr("failed to get bandwidth groups", err)
		}

		return nil, BandwidthGroupGetResult{
			Count:  len(groups),
			Output: formatBandwidthGroups(groups),
		}, nil
	}
}

// BandwidthGroupGetTool returns the MCP tool definition for transmission_bandwidth_group_get.
func BandwidthGroupGetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_bandwidth_group_get",
		Description: "Get bandwidth group configurations",
	}
}

func formatBandwidthGroups(groups []transmission.BandwidthGroup) string {
	if len(groups) == 0 {
		return "No bandwidth groups configured."
	}

	var bld strings.Builder

	fmt.Fprintf(&bld, "Bandwidth Groups (%d):\n", len(groups))

	for _, grp := range groups {
		fmt.Fprintf(&bld, "\n  %s:\n", grp.Name)

		if derefBool(grp.SpeedLimitDownEnabled) {
			fmt.Fprintf(&bld, "    DL Limit: %d KB/s\n", derefInt64(grp.SpeedLimitDown))
		} else {
			bld.WriteString("    DL Limit: unlimited\n")
		}

		if derefBool(grp.SpeedLimitUpEnabled) {
			fmt.Fprintf(&bld, "    UL Limit: %d KB/s\n", derefInt64(grp.SpeedLimitUp))
		} else {
			bld.WriteString("    UL Limit: unlimited\n")
		}

		fmt.Fprintf(&bld, "    Honors Session Limits: %v\n", derefBool(grp.HonorsSessionLimits))
	}

	return bld.String()
}

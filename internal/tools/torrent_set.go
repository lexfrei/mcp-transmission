package tools

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrNoTorrentChanges is returned when no torrent parameters are provided.
var ErrNoTorrentChanges = errors.New("at least one torrent parameter must be provided")

// TorrentSetParams defines the parameters for the transmission_torrent_set tool.
type TorrentSetParams struct {
	IDs             []int64  `json:"ids"                       jsonschema:"Torrent IDs to modify"`
	DownloadLimit   *int64   `json:"downloadLimit,omitempty"   jsonschema:"Download speed limit in KB/s"`
	DownloadLimited *bool    `json:"downloadLimited,omitempty" jsonschema:"Enable download speed limit"`
	UploadLimit     *int64   `json:"uploadLimit,omitempty"     jsonschema:"Upload speed limit in KB/s"`
	UploadLimited   *bool    `json:"uploadLimited,omitempty"   jsonschema:"Enable upload speed limit"`
	Labels          []string `json:"labels,omitempty"          jsonschema:"Labels/tags for the torrents"`
	SeedRatioLimit  *float64 `json:"seedRatioLimit,omitempty"  jsonschema:"Seed ratio limit"`
	QueuePosition   *int     `json:"queuePosition,omitempty"   jsonschema:"Queue position"`
}

// TorrentSetResult is the output of the transmission_torrent_set tool.
type TorrentSetResult struct {
	Message string `json:"message"`
}

// NewTorrentSetHandler creates a handler for the transmission_torrent_set tool.
func NewTorrentSetHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentSetParams, TorrentSetResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentSetParams,
	) (*mcp.CallToolResult, TorrentSetResult, error) {
		if len(params.IDs) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentSetResult{},
				validationErr(ErrIDsRequired)
		}

		if !hasTorrentChanges(&params) {
			return &mcp.CallToolResult{IsError: true}, TorrentSetResult{},
				validationErr(ErrNoTorrentChanges)
		}

		args := buildTorrentSetArgs(&params)

		err := client.TorrentSet(ctx, params.IDs, args)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentSetResult{},
				transmissionErr("failed to modify torrents", err)
		}

		msg := fmt.Sprintf("Modified %d torrent(s)", len(params.IDs))

		return nil, TorrentSetResult{Message: msg}, nil
	}
}

// TorrentSetTool returns the MCP tool definition for transmission_torrent_set.
func TorrentSetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_set",
		Description: "Modify properties of one or more torrents (speed limits, labels, seed ratio, etc.)",
	}
}

func hasTorrentChanges(params *TorrentSetParams) bool {
	return params.DownloadLimit != nil ||
		params.DownloadLimited != nil ||
		params.UploadLimit != nil ||
		params.UploadLimited != nil ||
		len(params.Labels) > 0 ||
		params.SeedRatioLimit != nil ||
		params.QueuePosition != nil
}

func buildTorrentSetArgs(params *TorrentSetParams) *transmission.TorrentSetArgs {
	args := &transmission.TorrentSetArgs{
		DownloadLimit:   params.DownloadLimit,
		DownloadLimited: params.DownloadLimited,
		UploadLimit:     params.UploadLimit,
		UploadLimited:   params.UploadLimited,
		SeedRatioLimit:  params.SeedRatioLimit,
		QueuePosition:   params.QueuePosition,
	}

	if len(params.Labels) > 0 {
		args.Labels = params.Labels
	}

	return args
}

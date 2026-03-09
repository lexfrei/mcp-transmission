package tools

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrFilenameOrMetainfoRequired is returned when neither filename nor metainfo is provided.
var ErrFilenameOrMetainfoRequired = errors.New("either filename (magnet/URL) or metainfo (base64 .torrent) is required")

// ErrFilenameAndMetainfoConflict is returned when both filename and metainfo are provided.
var ErrFilenameAndMetainfoConflict = errors.New("filename and metainfo are mutually exclusive, provide only one")

// TorrentAddParams defines the parameters for the transmission_torrent_add tool.
type TorrentAddParams struct {
	Filename    string   `json:"filename,omitempty"    jsonschema:"Magnet link, URL, or torrent file path"`
	Metainfo    string   `json:"metainfo,omitempty"    jsonschema:"Base64-encoded .torrent file content"`
	DownloadDir string   `json:"downloadDir,omitempty" jsonschema:"Download directory path"`
	Paused      *bool    `json:"paused,omitempty"      jsonschema:"Start torrent paused"`
	Labels      []string `json:"labels,omitempty"      jsonschema:"Labels/tags for the torrent"`
}

// TorrentAddResult is the output of the transmission_torrent_add tool.
type TorrentAddResult struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Duplicate bool   `json:"duplicate"`
	Message   string `json:"message"`
}

// NewTorrentAddHandler creates a handler for the transmission_torrent_add tool.
func NewTorrentAddHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentAddParams, TorrentAddResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentAddParams,
	) (*mcp.CallToolResult, TorrentAddResult, error) {
		if params.Filename == "" && params.Metainfo == "" {
			return &mcp.CallToolResult{IsError: true}, TorrentAddResult{},
				validationErr(ErrFilenameOrMetainfoRequired)
		}

		if params.Filename != "" && params.Metainfo != "" {
			return &mcp.CallToolResult{IsError: true}, TorrentAddResult{},
				validationErr(ErrFilenameAndMetainfoConflict)
		}

		args := buildTorrentAddArgs(&params)

		resp, err := client.TorrentAdd(ctx, args)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentAddResult{},
				transmissionErr("failed to add torrent", err)
		}

		return nil, buildTorrentAddResult(resp), nil
	}
}

// TorrentAddTool returns the MCP tool definition for transmission_torrent_add.
func TorrentAddTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_add",
		Description: "Add a new torrent by magnet link, URL, or base64-encoded .torrent file content",
	}
}

func buildTorrentAddArgs(params *TorrentAddParams) *transmission.TorrentAddArgs {
	args := &transmission.TorrentAddArgs{}

	if params.Filename != "" {
		args.Filename = &params.Filename
	}

	if params.Metainfo != "" {
		args.Metainfo = &params.Metainfo
	}

	if params.DownloadDir != "" {
		args.DownloadDir = &params.DownloadDir
	}

	if params.Paused != nil {
		args.Paused = params.Paused
	}

	if len(params.Labels) > 0 {
		args.Labels = params.Labels
	}

	return args
}

func buildTorrentAddResult(resp *transmission.TorrentAddResult) TorrentAddResult {
	if resp == nil {
		return TorrentAddResult{Message: "Torrent added (no details returned)"}
	}

	if resp.TorrentAdded != nil {
		return TorrentAddResult{
			ID:      resp.TorrentAdded.ID,
			Name:    resp.TorrentAdded.Name,
			Hash:    resp.TorrentAdded.HashString,
			Message: fmt.Sprintf("Added torrent: %s (ID: %d)", resp.TorrentAdded.Name, resp.TorrentAdded.ID),
		}
	}

	if resp.TorrentDuplicate != nil {
		return TorrentAddResult{
			ID:        resp.TorrentDuplicate.ID,
			Name:      resp.TorrentDuplicate.Name,
			Hash:      resp.TorrentDuplicate.HashString,
			Duplicate: true,
			Message:   fmt.Sprintf("Torrent already exists: %s (ID: %d)", resp.TorrentDuplicate.Name, resp.TorrentDuplicate.ID),
		}
	}

	return TorrentAddResult{Message: "Torrent added (no details returned)"}
}

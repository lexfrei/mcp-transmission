package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// torrentListFields returns the fields requested for listing torrents.
func torrentListFields() []string {
	return []string{
		"id", "name", "status", "percentDone", "totalSize",
		"rateDownload", "rateUpload", "eta", "uploadedEver",
		"downloadedEver", "labels", "errorString", "error",
		"addedDate", "queuePosition",
	}
}

// TorrentListParams defines the parameters for the transmission_torrent_list tool.
type TorrentListParams struct{}

// TorrentListResult is the output of the transmission_torrent_list tool.
type TorrentListResult struct {
	Count   int    `json:"count"`
	Output  string `json:"output"`
	Summary string `json:"summary"`
}

// NewTorrentListHandler creates a handler for the transmission_torrent_list tool.
func NewTorrentListHandler(client transmission.Client) mcp.ToolHandlerFor[TorrentListParams, TorrentListResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ TorrentListParams,
	) (*mcp.CallToolResult, TorrentListResult, error) {
		result, err := client.TorrentGet(ctx, torrentListFields(), nil)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentListResult{},
				transmissionErr("failed to list torrents", err)
		}

		output := formatTorrentList(result.Torrents)

		return nil, TorrentListResult{
			Count:   len(result.Torrents),
			Output:  output,
			Summary: formatListSummary(result.Torrents),
		}, nil
	}
}

// TorrentListTool returns the MCP tool definition for transmission_torrent_list.
func TorrentListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_list",
		Description: "List all torrents with their status, progress, speed, and basic info",
	}
}

func formatTorrentList(torrents []transmission.Torrent) string {
	if len(torrents) == 0 {
		return "No torrents found."
	}

	var bld strings.Builder

	for idx := range torrents {
		formatTorrentSummaryLine(&bld, &torrents[idx])
	}

	return bld.String()
}

func formatTorrentSummaryLine(bld *strings.Builder, tor *transmission.Torrent) {
	name := deref(tor.Name, "Unknown")
	tid := derefInt64(tor.ID)
	status := derefStatus(tor.Status)
	pct := derefFloat64(tor.PercentDone)
	rateDown := derefInt64(tor.RateDownload)
	rateUp := derefInt64(tor.RateUpload)

	fmt.Fprintf(bld, "[%d] %s\n", tid, name)
	fmt.Fprintf(bld, "    Status: %s | Progress: %s", status, formatPercent(pct))

	if rateDown > 0 || rateUp > 0 {
		fmt.Fprintf(bld, " | DL: %s | UL: %s", formatSpeed(rateDown), formatSpeed(rateUp))
	}

	bld.WriteString("\n")
}

func formatListSummary(torrents []transmission.Torrent) string {
	var downloading, seeding, stopped, checking int

	for idx := range torrents {
		if torrents[idx].Status == nil {
			continue
		}

		switch *torrents[idx].Status {
		case transmission.TorrentStatusDownload, transmission.TorrentStatusDownloadWait:
			downloading++
		case transmission.TorrentStatusSeed, transmission.TorrentStatusSeedWait:
			seeding++
		case transmission.TorrentStatusStopped:
			stopped++
		case transmission.TorrentStatusCheck, transmission.TorrentStatusCheckWait:
			checking++
		}
	}

	return fmt.Sprintf(
		"Total: %d | Downloading: %d | Seeding: %d | Stopped: %d | Checking: %d",
		len(torrents), downloading, seeding, stopped, checking,
	)
}

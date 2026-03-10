package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SessionStatsParams defines the parameters for the transmission_session_stats tool.
type SessionStatsParams struct{}

// SessionStatsResult is the output of the transmission_session_stats tool.
type SessionStatsResult struct {
	ActiveTorrents int    `json:"activeTorrents"`
	PausedTorrents int    `json:"pausedTorrents"`
	TotalTorrents  int    `json:"totalTorrents"`
	DownloadSpeed  string `json:"downloadSpeed"`
	UploadSpeed    string `json:"uploadSpeed"`
	Output         string `json:"output"`
}

// NewSessionStatsHandler creates a handler for the transmission_session_stats tool.
func NewSessionStatsHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[SessionStatsParams, SessionStatsResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ SessionStatsParams,
	) (*mcp.CallToolResult, SessionStatsResult, error) {
		stats, err := client.SessionStats(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SessionStatsResult{},
				transmissionErr("failed to get session stats", err)
		}

		if stats == nil {
			return nil, SessionStatsResult{Output: "No stats data returned"}, nil
		}

		result := SessionStatsResult{
			ActiveTorrents: stats.ActiveTorrentCount,
			PausedTorrents: stats.PausedTorrentCount,
			TotalTorrents:  stats.TorrentCount,
			DownloadSpeed:  formatSpeed(stats.DownloadSpeed),
			UploadSpeed:    formatSpeed(stats.UploadSpeed),
			Output:         formatSessionStats(stats),
		}

		return nil, result, nil
	}
}

// SessionStatsTool returns the MCP tool definition for transmission_session_stats.
func SessionStatsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_session_stats",
		Description: "Get Transmission session statistics: speeds, torrent counts, and cumulative transfer totals",
	}
}

func formatSessionStats(stats *transmission.SessionStats) string {
	var bld strings.Builder

	fmt.Fprintf(&bld, "Torrents: %d total (%d active, %d paused)\n",
		stats.TorrentCount, stats.ActiveTorrentCount, stats.PausedTorrentCount)
	fmt.Fprintf(&bld, "Speed: DL %s | UL %s\n",
		formatSpeed(stats.DownloadSpeed), formatSpeed(stats.UploadSpeed))

	bld.WriteString("\nCurrent Session:\n")
	fmt.Fprintf(&bld, "  Downloaded: %s\n", formatBytes(stats.CurrentStats.DownloadedBytes))
	fmt.Fprintf(&bld, "  Uploaded: %s\n", formatBytes(stats.CurrentStats.UploadedBytes))
	fmt.Fprintf(&bld, "  Files Added: %d\n", stats.CurrentStats.FilesAdded)

	bld.WriteString("\nCumulative:\n")
	fmt.Fprintf(&bld, "  Downloaded: %s\n", formatBytes(stats.CumulativeStats.DownloadedBytes))
	fmt.Fprintf(&bld, "  Uploaded: %s\n", formatBytes(stats.CumulativeStats.UploadedBytes))
	fmt.Fprintf(&bld, "  Files Added: %d\n", stats.CumulativeStats.FilesAdded)
	fmt.Fprintf(&bld, "  Sessions: %d\n", stats.CumulativeStats.SessionCount)

	return bld.String()
}

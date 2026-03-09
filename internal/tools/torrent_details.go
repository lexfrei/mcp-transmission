package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrPositiveIDRequired is returned when a positive torrent ID is needed.
var ErrPositiveIDRequired = errors.New("torrent ID must be positive")

// torrentDetailFields returns all fields requested for detailed torrent info.
func torrentDetailFields() []string {
	return []string{
		"id", "name", "status", "hashString", "totalSize", "sizeWhenDone",
		"leftUntilDone", "percentDone", "percentComplete", "rateDownload",
		"rateUpload", "uploadedEver", "downloadedEver", "eta", "etaIdle",
		"error", "errorString", "peersConnected", "peersSendingToUs",
		"peersGettingFromUs", "addedDate", "doneDate", "downloadDir",
		"comment", "creator", "isPrivate", "labels", "magnetLink",
		"fileCount", "files", "trackers", "trackerStats", "queuePosition",
		"bandwidthPriority", "downloadLimit", "downloadLimited",
		"uploadLimit", "uploadLimited", "seedRatioLimit", "seedRatioMode",
		"activityDate", "secondsDownloading", "secondsSeeding",
		"desiredAvailable", "haveValid", "isFinished", "isStalled",
	}
}

// TorrentDetailsParams defines the parameters for the transmission_torrent_details tool.
type TorrentDetailsParams struct {
	ID int64 `json:"id" jsonschema:"Torrent ID"`
}

// TorrentDetailsResult is the output of the transmission_torrent_details tool.
type TorrentDetailsResult struct {
	Output string `json:"output"`
}

// NewTorrentDetailsHandler creates a handler for the transmission_torrent_details tool.
func NewTorrentDetailsHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[TorrentDetailsParams, TorrentDetailsResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TorrentDetailsParams,
	) (*mcp.CallToolResult, TorrentDetailsResult, error) {
		if params.ID <= 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentDetailsResult{},
				validationErr(ErrPositiveIDRequired)
		}

		result, err := client.TorrentGet(ctx, torrentDetailFields(), []int64{params.ID})
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, TorrentDetailsResult{},
				transmissionErr("failed to get torrent details", err)
		}

		if len(result.Torrents) == 0 {
			return &mcp.CallToolResult{IsError: true}, TorrentDetailsResult{},
				validationErr(fmt.Errorf("torrent %d not found", params.ID)) //nolint:err113 // dynamic ID in message.
		}

		output := formatTorrentDetails(&result.Torrents[0])

		return nil, TorrentDetailsResult{Output: output}, nil
	}
}

// TorrentDetailsTool returns the MCP tool definition for transmission_torrent_details.
func TorrentDetailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_torrent_details",
		Description: "Get detailed information about a specific torrent including files, trackers, and peers",
	}
}

func formatTorrentDetails(tor *transmission.Torrent) string {
	var bld strings.Builder

	fmt.Fprintf(&bld, "Name: %s\n", deref(tor.Name, "Unknown"))
	fmt.Fprintf(&bld, "ID: %d\n", derefInt64(tor.ID))
	fmt.Fprintf(&bld, "Hash: %s\n", deref(tor.HashString, ""))
	fmt.Fprintf(&bld, "Status: %s\n", derefStatus(tor.Status))
	fmt.Fprintf(&bld, "Progress: %s\n", formatPercent(derefFloat64(tor.PercentDone)))
	fmt.Fprintf(&bld, "Size: %s\n", formatBytes(derefInt64(tor.TotalSize)))
	fmt.Fprintf(&bld, "Downloaded: %s\n", formatBytes(derefInt64(tor.DownloadedEver)))
	fmt.Fprintf(&bld, "Uploaded: %s\n", formatBytes(derefInt64(tor.UploadedEver)))
	fmt.Fprintf(&bld, "DL Speed: %s\n", formatSpeed(derefInt64(tor.RateDownload)))
	fmt.Fprintf(&bld, "UL Speed: %s\n", formatSpeed(derefInt64(tor.RateUpload)))

	formatDetailsETA(&bld, tor)
	formatDetailsPeers(&bld, tor)
	formatDetailsLocation(&bld, tor)
	formatDetailsFiles(&bld, tor)
	formatDetailsTrackers(&bld, tor)
	formatDetailsLabels(&bld, tor)

	return bld.String()
}

func formatDetailsETA(bld *strings.Builder, tor *transmission.Torrent) {
	if tor.ETA != nil && *tor.ETA >= 0 {
		dur := time.Duration(derefInt64(tor.ETA)) * time.Second
		fmt.Fprintf(bld, "ETA: %s\n", dur)
	}
}

func formatDetailsPeers(bld *strings.Builder, tor *transmission.Torrent) {
	fmt.Fprintf(bld, "Peers: %d connected", derefInt(tor.PeersConnected))

	if tor.PeersSendingToUs != nil {
		fmt.Fprintf(bld, ", %d sending", *tor.PeersSendingToUs)
	}

	if tor.PeersGettingFromUs != nil {
		fmt.Fprintf(bld, ", %d receiving", *tor.PeersGettingFromUs)
	}

	bld.WriteString("\n")
}

func formatDetailsLocation(bld *strings.Builder, tor *transmission.Torrent) {
	fmt.Fprintf(bld, "Location: %s\n", deref(tor.DownloadDir, ""))

	if tor.Comment != nil && *tor.Comment != "" {
		fmt.Fprintf(bld, "Comment: %s\n", *tor.Comment)
	}

	if tor.IsPrivate != nil {
		fmt.Fprintf(bld, "Private: %v\n", *tor.IsPrivate)
	}
}

func formatDetailsFiles(bld *strings.Builder, tor *transmission.Torrent) {
	if len(tor.Files) == 0 {
		return
	}

	fmt.Fprintf(bld, "\nFiles (%d):\n", len(tor.Files))

	for idx, file := range tor.Files {
		pct := float64(0)
		if file.Length > 0 {
			pct = float64(file.BytesCompleted) / float64(file.Length)
		}

		fmt.Fprintf(bld, "  [%d] %s (%s, %s)\n", idx, file.Name, formatBytes(file.Length), formatPercent(pct))
	}
}

func formatDetailsTrackers(bld *strings.Builder, tor *transmission.Torrent) {
	if len(tor.TrackerStats) == 0 {
		return
	}

	bld.WriteString("\nTrackers:\n")

	for idx := range tor.TrackerStats {
		fmt.Fprintf(bld, "  %s (Seeds: %d, Leechers: %d, Last: %s)\n",
			tor.TrackerStats[idx].Host, tor.TrackerStats[idx].SeederCount,
			tor.TrackerStats[idx].LeecherCount, tor.TrackerStats[idx].LastAnnounceResult)
	}
}

func formatDetailsLabels(bld *strings.Builder, tor *transmission.Torrent) {
	if len(tor.Labels) == 0 {
		return
	}

	fmt.Fprintf(bld, "Labels: %s\n", strings.Join(tor.Labels, ", "))
}

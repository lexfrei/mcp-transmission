package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SessionGetParams defines the parameters for the transmission_session_get tool.
type SessionGetParams struct{}

// SessionGetResult is the output of the transmission_session_get tool.
type SessionGetResult struct {
	Output string `json:"output"`
}

// NewSessionGetHandler creates a handler for the transmission_session_get tool.
func NewSessionGetHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[SessionGetParams, SessionGetResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ SessionGetParams,
	) (*mcp.CallToolResult, SessionGetResult, error) {
		session, err := client.SessionGet(ctx, nil)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SessionGetResult{},
				transmissionErr("failed to get session config", err)
		}

		return nil, SessionGetResult{Output: formatSession(session)}, nil
	}
}

// SessionGetTool returns the MCP tool definition for transmission_session_get.
func SessionGetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_session_get",
		Description: "Get Transmission session configuration: directories, speed limits, encryption, peer settings, etc.",
	}
}

func formatSession(ses *transmission.Session) string {
	var bld strings.Builder

	bld.WriteString("Transmission Session Configuration:\n\n")

	formatSessionVersion(&bld, ses)
	formatSessionDirs(&bld, ses)
	formatSessionSpeed(&bld, ses)
	formatSessionPeers(&bld, ses)
	formatSessionQueue(&bld, ses)

	return bld.String()
}

func formatSessionVersion(bld *strings.Builder, ses *transmission.Session) {
	if ses.Version != nil {
		fmt.Fprintf(bld, "Version: %s\n", *ses.Version)
	}

	if ses.RPCVersionSemver != nil {
		fmt.Fprintf(bld, "RPC Version: %s\n", *ses.RPCVersionSemver)
	}
}

func formatSessionDirs(bld *strings.Builder, ses *transmission.Session) {
	bld.WriteString("\nDirectories:\n")
	fmt.Fprintf(bld, "  Download: %s\n", deref(ses.DownloadDir, ""))

	if derefBool(ses.IncompleteDirEnabled) {
		fmt.Fprintf(bld, "  Incomplete: %s\n", deref(ses.IncompleteDir, ""))
	}
}

func formatSessionSpeed(bld *strings.Builder, ses *transmission.Session) {
	bld.WriteString("\nSpeed Limits:\n")

	if derefBool(ses.SpeedLimitDownEnabled) {
		fmt.Fprintf(bld, "  Download: %d KB/s\n", derefInt64(ses.SpeedLimitDown))
	} else {
		bld.WriteString("  Download: unlimited\n")
	}

	if derefBool(ses.SpeedLimitUpEnabled) {
		fmt.Fprintf(bld, "  Upload: %d KB/s\n", derefInt64(ses.SpeedLimitUp))
	} else {
		bld.WriteString("  Upload: unlimited\n")
	}

	if derefBool(ses.AltSpeedEnabled) {
		fmt.Fprintf(bld, "  Alt Speed (turtle): DL %d / UL %d KB/s\n",
			derefInt64(ses.AltSpeedDown), derefInt64(ses.AltSpeedUp))
	}
}

func formatSessionPeers(bld *strings.Builder, ses *transmission.Session) {
	bld.WriteString("\nPeer Settings:\n")
	fmt.Fprintf(bld, "  Global Limit: %d\n", derefInt(ses.PeerLimitGlobal))
	fmt.Fprintf(bld, "  Per Torrent: %d\n", derefInt(ses.PeerLimitPerTorrent))
	fmt.Fprintf(bld, "  Port: %d\n", derefInt(ses.PeerPort))
	fmt.Fprintf(bld, "  DHT: %v | PEX: %v | LPD: %v\n",
		derefBool(ses.DHTEnabled), derefBool(ses.PEXEnabled), derefBool(ses.LPDEnabled))

	if ses.Encryption != nil {
		fmt.Fprintf(bld, "  Encryption: %s\n", *ses.Encryption)
	}
}

func formatSessionQueue(bld *strings.Builder, ses *transmission.Session) {
	bld.WriteString("\nQueue:\n")
	fmt.Fprintf(bld, "  Download Queue: %v (size: %d)\n",
		derefBool(ses.DownloadQueueEnabled), derefInt(ses.DownloadQueueSize))
	fmt.Fprintf(bld, "  Seed Queue: %v (size: %d)\n",
		derefBool(ses.SeedQueueEnabled), derefInt(ses.SeedQueueSize))
}

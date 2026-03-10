package tools

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/go-transmission/api/transmission"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrNoSessionChanges is returned when no session parameters are provided.
var ErrNoSessionChanges = errors.New("at least one session parameter must be provided")

// ErrEmptyDownloadDir is returned when downloadDir is an empty string.
var ErrEmptyDownloadDir = errors.New("downloadDir must not be empty")

// SessionSetParams defines the parameters for the transmission_session_set tool.
type SessionSetParams struct {
	SpeedLimitDown        *int64  `json:"speedLimitDown,omitempty"        jsonschema:"Download speed limit in KB/s"`
	SpeedLimitDownEnabled *bool   `json:"speedLimitDownEnabled,omitempty" jsonschema:"Enable download speed limit"`
	SpeedLimitUp          *int64  `json:"speedLimitUp,omitempty"          jsonschema:"Upload speed limit in KB/s"`
	SpeedLimitUpEnabled   *bool   `json:"speedLimitUpEnabled,omitempty"   jsonschema:"Enable upload speed limit"`
	AltSpeedDown          *int64  `json:"altSpeedDown,omitempty"          jsonschema:"Alternative (turtle) download limit in KB/s"`
	AltSpeedUp            *int64  `json:"altSpeedUp,omitempty"            jsonschema:"Alternative (turtle) upload limit in KB/s"`
	AltSpeedEnabled       *bool   `json:"altSpeedEnabled,omitempty"       jsonschema:"Enable alternative speed mode"`
	DownloadDir           *string `json:"downloadDir,omitempty"           jsonschema:"Default download directory"`
	PeerLimitGlobal       *int    `json:"peerLimitGlobal,omitempty"       jsonschema:"Global peer limit"`
	PeerLimitPerTorrent   *int    `json:"peerLimitPerTorrent,omitempty"   jsonschema:"Per-torrent peer limit"`
}

// SessionSetResult is the output of the transmission_session_set tool.
type SessionSetResult struct {
	Message string `json:"message"`
}

// NewSessionSetHandler creates a handler for the transmission_session_set tool.
func NewSessionSetHandler(
	client transmission.Client,
) mcp.ToolHandlerFor[SessionSetParams, SessionSetResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params SessionSetParams,
	) (*mcp.CallToolResult, SessionSetResult, error) {
		if !hasSessionChanges(&params) {
			return &mcp.CallToolResult{IsError: true}, SessionSetResult{},
				validationErr(ErrNoSessionChanges)
		}

		limErr := validateSessionLimits(&params)
		if limErr != nil {
			return &mcp.CallToolResult{IsError: true}, SessionSetResult{}, limErr
		}

		args := buildSessionSetArgs(&params)

		err := client.SessionSet(ctx, args)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SessionSetResult{},
				transmissionErr("failed to update session", err)
		}

		return nil, SessionSetResult{
			Message: describeSessionChanges(&params),
		}, nil
	}
}

// SessionSetTool returns the MCP tool definition for transmission_session_set.
func SessionSetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "transmission_session_set",
		Description: "Modify Transmission session settings: speed limits, alt speed (turtle mode), download directory, peer limits",
	}
}

func hasSessionChanges(params *SessionSetParams) bool {
	return params.SpeedLimitDown != nil ||
		params.SpeedLimitDownEnabled != nil ||
		params.SpeedLimitUp != nil ||
		params.SpeedLimitUpEnabled != nil ||
		params.AltSpeedDown != nil ||
		params.AltSpeedUp != nil ||
		params.AltSpeedEnabled != nil ||
		params.DownloadDir != nil ||
		params.PeerLimitGlobal != nil ||
		params.PeerLimitPerTorrent != nil
}

func validateSessionLimits(params *SessionSetParams) error {
	speedErr := validateSessionSpeedLimits(params)
	if speedErr != nil {
		return speedErr
	}

	if params.PeerLimitGlobal != nil && *params.PeerLimitGlobal < 0 {
		return validationErr(ErrNegativeLimit)
	}

	if params.PeerLimitPerTorrent != nil && *params.PeerLimitPerTorrent < 0 {
		return validationErr(ErrNegativeLimit)
	}

	if params.DownloadDir != nil && *params.DownloadDir == "" {
		return validationErr(ErrEmptyDownloadDir)
	}

	if params.DownloadDir != nil && !strings.HasPrefix(*params.DownloadDir, "/") {
		return validationErr(ErrAbsolutePathRequired)
	}

	return nil
}

func validateSessionSpeedLimits(params *SessionSetParams) error {
	if params.SpeedLimitDown != nil && *params.SpeedLimitDown < 0 {
		return validationErr(ErrNegativeLimit)
	}

	if params.SpeedLimitUp != nil && *params.SpeedLimitUp < 0 {
		return validationErr(ErrNegativeLimit)
	}

	if params.AltSpeedDown != nil && *params.AltSpeedDown < 0 {
		return validationErr(ErrNegativeLimit)
	}

	if params.AltSpeedUp != nil && *params.AltSpeedUp < 0 {
		return validationErr(ErrNegativeLimit)
	}

	return nil
}

func buildSessionSetArgs(params *SessionSetParams) *transmission.SessionSetArgs {
	return &transmission.SessionSetArgs{
		SpeedLimitDown:        params.SpeedLimitDown,
		SpeedLimitDownEnabled: params.SpeedLimitDownEnabled,
		SpeedLimitUp:          params.SpeedLimitUp,
		SpeedLimitUpEnabled:   params.SpeedLimitUpEnabled,
		AltSpeedDown:          params.AltSpeedDown,
		AltSpeedUp:            params.AltSpeedUp,
		AltSpeedEnabled:       params.AltSpeedEnabled,
		DownloadDir:           params.DownloadDir,
		PeerLimitGlobal:       params.PeerLimitGlobal,
		PeerLimitPerTorrent:   params.PeerLimitPerTorrent,
	}
}

func describeSessionChanges(params *SessionSetParams) string {
	var changes []string

	if params.SpeedLimitDown != nil || params.SpeedLimitDownEnabled != nil {
		changes = append(changes, "download speed limit")
	}

	if params.SpeedLimitUp != nil || params.SpeedLimitUpEnabled != nil {
		changes = append(changes, "upload speed limit")
	}

	if params.AltSpeedDown != nil || params.AltSpeedUp != nil || params.AltSpeedEnabled != nil {
		changes = append(changes, "alternative speed settings")
	}

	if params.DownloadDir != nil {
		changes = append(changes, "download directory")
	}

	if params.PeerLimitGlobal != nil || params.PeerLimitPerTorrent != nil {
		changes = append(changes, "peer limits")
	}

	return "Updated session: " + strings.Join(changes, ", ")
}

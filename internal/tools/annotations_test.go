package tools_test

import (
	"testing"

	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// allTools returns every tool definition for table-driven tests.
func allTools() []*mcp.Tool {
	return []*mcp.Tool{
		tools.TorrentListTool(),
		tools.TorrentAddTool(),
		tools.TorrentRemoveTool(),
		tools.TorrentStartTool(),
		tools.TorrentStopTool(),
		tools.TorrentVerifyTool(),
		tools.TorrentReannounceTool(),
		tools.TorrentDetailsTool(),
		tools.TorrentSetTool(),
		tools.TorrentMoveTool(),
		tools.SessionStatsTool(),
		tools.SessionGetTool(),
		tools.SessionSetTool(),
		tools.FreeSpaceTool(),
		tools.PortTestTool(),
		tools.BlocklistUpdateTool(),
		tools.QueueMoveTool(),
		tools.BandwidthGroupGetTool(),
	}
}

func TestAllTools_HaveAnnotations(t *testing.T) {
	for _, tool := range allTools() {
		if tool.Annotations == nil {
			t.Errorf("tool %s: missing Annotations", tool.Name)
		}
	}
}

func TestAllTools_HaveTitle(t *testing.T) {
	for _, tool := range allTools() {
		if tool.Annotations == nil {
			t.Errorf("tool %s: missing Annotations", tool.Name)

			continue
		}

		if tool.Annotations.Title == "" {
			t.Errorf("tool %s: missing Title in Annotations", tool.Name)
		}
	}
}

func TestReadOnlyTools_HaveAnnotations(t *testing.T) {
	readOnlyTools := []*mcp.Tool{
		tools.TorrentListTool(),
		tools.TorrentDetailsTool(),
		tools.SessionStatsTool(),
		tools.SessionGetTool(),
		tools.FreeSpaceTool(),
		tools.PortTestTool(),
		tools.BandwidthGroupGetTool(),
	}

	for _, tool := range readOnlyTools {
		if tool.Annotations == nil {
			t.Errorf("tool %s: missing Annotations", tool.Name)

			continue
		}

		if !tool.Annotations.ReadOnlyHint {
			t.Errorf("tool %s: expected ReadOnlyHint=true", tool.Name)
		}
	}
}

func TestIdempotentWriteTools_HaveAnnotations(t *testing.T) {
	idempotentTools := []*mcp.Tool{
		tools.TorrentStartTool(),
		tools.TorrentStopTool(),
		tools.TorrentVerifyTool(),
		tools.TorrentReannounceTool(),
		tools.TorrentSetTool(),
		tools.TorrentMoveTool(),
		tools.SessionSetTool(),
		tools.BlocklistUpdateTool(),
		tools.QueueMoveTool(),
		tools.TorrentRemoveTool(),
	}

	for _, tool := range idempotentTools {
		if tool.Annotations == nil {
			t.Errorf("tool %s: missing Annotations", tool.Name)

			continue
		}

		if !tool.Annotations.IdempotentHint {
			t.Errorf("tool %s: expected IdempotentHint=true", tool.Name)
		}
	}
}

func TestNonDestructiveWriteTools_HaveAnnotations(t *testing.T) {
	nonDestructiveTools := []*mcp.Tool{
		tools.TorrentAddTool(),
		tools.TorrentStartTool(),
		tools.TorrentStopTool(),
		tools.TorrentVerifyTool(),
		tools.TorrentReannounceTool(),
		tools.TorrentSetTool(),
		tools.TorrentMoveTool(),
		tools.SessionSetTool(),
		tools.BlocklistUpdateTool(),
		tools.QueueMoveTool(),
	}

	for _, tool := range nonDestructiveTools {
		if tool.Annotations == nil {
			t.Errorf("tool %s: missing Annotations", tool.Name)

			continue
		}

		if tool.Annotations.DestructiveHint == nil {
			t.Errorf("tool %s: DestructiveHint must be explicitly set to false (nil defaults to true)", tool.Name)

			continue
		}

		if *tool.Annotations.DestructiveHint {
			t.Errorf("tool %s: expected DestructiveHint=false", tool.Name)
		}
	}
}

func TestDestructiveTools_HaveAnnotations(t *testing.T) {
	tool := tools.TorrentRemoveTool()

	if tool.Annotations == nil {
		t.Fatal("torrent_remove: missing Annotations")
	}

	// DestructiveHint defaults to true when nil, but we want it explicit.
	if tool.Annotations.DestructiveHint != nil && !*tool.Annotations.DestructiveHint {
		t.Error("torrent_remove: expected DestructiveHint=true (or nil for default)")
	}
}

func TestTorrentAddTool_NotIdempotent(t *testing.T) {
	tool := tools.TorrentAddTool()

	if tool.Annotations == nil {
		t.Fatal("torrent_add: missing Annotations")
	}

	if tool.Annotations.IdempotentHint {
		t.Error("torrent_add: expected IdempotentHint=false (adding a torrent is not idempotent)")
	}
}

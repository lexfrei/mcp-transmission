// Package testutil provides shared test helpers for the mcp-transmission project.
package testutil

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
)

// NoopClient implements transmission.Client with no-op methods for testing.
type NoopClient struct{}

func (n *NoopClient) TorrentStart(_ context.Context, _ []int64) error      { return nil }
func (n *NoopClient) TorrentStartNow(_ context.Context, _ []int64) error   { return nil }
func (n *NoopClient) TorrentStop(_ context.Context, _ []int64) error       { return nil }
func (n *NoopClient) TorrentVerify(_ context.Context, _ []int64) error     { return nil }
func (n *NoopClient) TorrentReannounce(_ context.Context, _ []int64) error { return nil }

func (n *NoopClient) TorrentGet(
	_ context.Context, _ []string, _ []int64,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *NoopClient) TorrentGetByHash(
	_ context.Context, _, _ []string,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *NoopClient) TorrentGetRecentlyActive(
	_ context.Context, _ []string,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *NoopClient) TorrentSet(_ context.Context, _ []int64, _ *transmission.TorrentSetArgs) error {
	return nil
}

func (n *NoopClient) TorrentAdd(
	_ context.Context, _ *transmission.TorrentAddArgs,
) (*transmission.TorrentAddResult, error) {
	return &transmission.TorrentAddResult{}, nil
}

func (n *NoopClient) TorrentRemove(_ context.Context, _ []int64, _ bool) error { return nil }

func (n *NoopClient) TorrentSetLocation(_ context.Context, _ []int64, _ string, _ bool) error {
	return nil
}

func (n *NoopClient) TorrentRenamePath(
	_ context.Context, _ int64, _, _ string,
) (*transmission.TorrentRenameResult, error) {
	return &transmission.TorrentRenameResult{}, nil
}

func (n *NoopClient) SessionGet(
	_ context.Context, _ []string,
) (*transmission.Session, error) {
	return &transmission.Session{}, nil
}

func (n *NoopClient) SessionSet(_ context.Context, _ *transmission.SessionSetArgs) error {
	return nil
}

func (n *NoopClient) SessionStats(_ context.Context) (*transmission.SessionStats, error) {
	return &transmission.SessionStats{}, nil
}

func (n *NoopClient) SessionClose(_ context.Context) error               { return nil }
func (n *NoopClient) QueueMoveTop(_ context.Context, _ []int64) error    { return nil }
func (n *NoopClient) QueueMoveUp(_ context.Context, _ []int64) error     { return nil }
func (n *NoopClient) QueueMoveDown(_ context.Context, _ []int64) error   { return nil }
func (n *NoopClient) QueueMoveBottom(_ context.Context, _ []int64) error { return nil }

func (n *NoopClient) BlocklistUpdate(_ context.Context) (int, error) { return 0, nil }
func (n *NoopClient) PortTest(_ context.Context) (bool, error)       { return false, nil }

func (n *NoopClient) FreeSpace(_ context.Context, _ string) (*transmission.FreeSpace, error) {
	return &transmission.FreeSpace{}, nil
}

func (n *NoopClient) GroupSet(_ context.Context, _ *transmission.BandwidthGroup) error { return nil }

func (n *NoopClient) GroupGet(
	_ context.Context, _ []string,
) ([]transmission.BandwidthGroup, error) {
	return nil, nil
}

func (n *NoopClient) Close() error { return nil }

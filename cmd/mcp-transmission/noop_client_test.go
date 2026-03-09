package main

import (
	"context"

	"github.com/lexfrei/go-transmission/api/transmission"
)

// noopClient implements transmission.Client for testing registration.
type noopClient struct{}

func (n *noopClient) TorrentStart(_ context.Context, _ []int64) error      { return nil }
func (n *noopClient) TorrentStartNow(_ context.Context, _ []int64) error   { return nil }
func (n *noopClient) TorrentStop(_ context.Context, _ []int64) error       { return nil }
func (n *noopClient) TorrentVerify(_ context.Context, _ []int64) error     { return nil }
func (n *noopClient) TorrentReannounce(_ context.Context, _ []int64) error { return nil }

func (n *noopClient) TorrentGet(
	_ context.Context, _ []string, _ []int64,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *noopClient) TorrentGetByHash(
	_ context.Context, _, _ []string,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *noopClient) TorrentGetRecentlyActive(
	_ context.Context, _ []string,
) (*transmission.TorrentGetResult, error) {
	return &transmission.TorrentGetResult{}, nil
}

func (n *noopClient) TorrentSet(_ context.Context, _ []int64, _ *transmission.TorrentSetArgs) error {
	return nil
}

func (n *noopClient) TorrentAdd(
	_ context.Context, _ *transmission.TorrentAddArgs,
) (*transmission.TorrentAddResult, error) {
	return &transmission.TorrentAddResult{}, nil
}

func (n *noopClient) TorrentRemove(_ context.Context, _ []int64, _ bool) error { return nil }

func (n *noopClient) TorrentSetLocation(_ context.Context, _ []int64, _ string, _ bool) error {
	return nil
}

func (n *noopClient) TorrentRenamePath(
	_ context.Context, _ int64, _, _ string,
) (*transmission.TorrentRenameResult, error) {
	return &transmission.TorrentRenameResult{}, nil
}

func (n *noopClient) SessionGet(
	_ context.Context, _ []string,
) (*transmission.Session, error) {
	return &transmission.Session{}, nil
}

func (n *noopClient) SessionSet(_ context.Context, _ *transmission.SessionSetArgs) error {
	return nil
}

func (n *noopClient) SessionStats(_ context.Context) (*transmission.SessionStats, error) {
	return &transmission.SessionStats{}, nil
}

func (n *noopClient) SessionClose(_ context.Context) error               { return nil }
func (n *noopClient) QueueMoveTop(_ context.Context, _ []int64) error    { return nil }
func (n *noopClient) QueueMoveUp(_ context.Context, _ []int64) error     { return nil }
func (n *noopClient) QueueMoveDown(_ context.Context, _ []int64) error   { return nil }
func (n *noopClient) QueueMoveBottom(_ context.Context, _ []int64) error { return nil }

func (n *noopClient) BlocklistUpdate(_ context.Context) (int, error) { return 0, nil }
func (n *noopClient) PortTest(_ context.Context) (bool, error)       { return false, nil }

func (n *noopClient) FreeSpace(_ context.Context, _ string) (*transmission.FreeSpace, error) {
	return &transmission.FreeSpace{}, nil
}

func (n *noopClient) GroupSet(_ context.Context, _ *transmission.BandwidthGroup) error { return nil }

func (n *noopClient) GroupGet(
	_ context.Context, _ []string,
) ([]transmission.BandwidthGroup, error) {
	return nil, nil
}

func (n *noopClient) Close() error { return nil }

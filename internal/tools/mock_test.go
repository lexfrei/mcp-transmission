package tools_test

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/go-transmission/api/transmission"
)

var errMock = errors.New("mock error")

// mockClient implements transmission.Client for testing.
type mockClient struct {
	err              error
	torrentGetResult *transmission.TorrentGetResult
	torrentAddResult *transmission.TorrentAddResult
	portTestResult   bool
	blocklistCount   int
	freeSpaceResult  *transmission.FreeSpace
	sessionResult    *transmission.Session
	sessionStats     *transmission.SessionStats
	bandwidthGroups  []transmission.BandwidthGroup
	torrentRenameRes *transmission.TorrentRenameResult
}

func newMockClient() *mockClient {
	return &mockClient{}
}

func (m *mockClient) TorrentStart(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) TorrentStartNow(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) TorrentStop(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) TorrentVerify(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) TorrentReannounce(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) TorrentGet(
	_ context.Context,
	_ []string,
	_ []int64,
) (*transmission.TorrentGetResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.torrentGetResult, nil
}

func (m *mockClient) TorrentGetByHash(
	_ context.Context,
	_, _ []string,
) (*transmission.TorrentGetResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.torrentGetResult, nil
}

func (m *mockClient) TorrentGetRecentlyActive(
	_ context.Context,
	_ []string,
) (*transmission.TorrentGetResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.torrentGetResult, nil
}

func (m *mockClient) TorrentSet(
	_ context.Context,
	_ []int64,
	_ *transmission.TorrentSetArgs,
) error {
	return m.err
}

func (m *mockClient) TorrentAdd(
	_ context.Context,
	_ *transmission.TorrentAddArgs,
) (*transmission.TorrentAddResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.torrentAddResult, nil
}

func (m *mockClient) TorrentRemove(_ context.Context, _ []int64, _ bool) error {
	return m.err
}

func (m *mockClient) TorrentSetLocation(
	_ context.Context,
	_ []int64,
	_ string,
	_ bool,
) error {
	return m.err
}

func (m *mockClient) TorrentRenamePath(
	_ context.Context,
	_ int64,
	_, _ string,
) (*transmission.TorrentRenameResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.torrentRenameRes, nil
}

func (m *mockClient) SessionGet(
	_ context.Context,
	_ []string,
) (*transmission.Session, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.sessionResult, nil
}

func (m *mockClient) SessionSet(
	_ context.Context,
	_ *transmission.SessionSetArgs,
) error {
	return m.err
}

func (m *mockClient) SessionStats(_ context.Context) (*transmission.SessionStats, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.sessionStats, nil
}

func (m *mockClient) SessionClose(_ context.Context) error {
	return m.err
}

func (m *mockClient) QueueMoveTop(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) QueueMoveUp(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) QueueMoveDown(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) QueueMoveBottom(_ context.Context, _ []int64) error {
	return m.err
}

func (m *mockClient) BlocklistUpdate(_ context.Context) (int, error) {
	if m.err != nil {
		return 0, m.err
	}

	return m.blocklistCount, nil
}

func (m *mockClient) PortTest(_ context.Context) (bool, error) {
	if m.err != nil {
		return false, m.err
	}

	return m.portTestResult, nil
}

func (m *mockClient) FreeSpace(
	_ context.Context,
	_ string,
) (*transmission.FreeSpace, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.freeSpaceResult, nil
}

func (m *mockClient) GroupSet(
	_ context.Context,
	_ *transmission.BandwidthGroup,
) error {
	return m.err
}

func (m *mockClient) GroupGet(
	_ context.Context,
	_ []string,
) ([]transmission.BandwidthGroup, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.bandwidthGroups, nil
}

func (m *mockClient) Close() error {
	return nil
}

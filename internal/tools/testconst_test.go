package tools_test

// Shared string fixtures reused across the tools test suite. Hoisted into a
// single place so the same literal is not repeated in multiple test files.
const (
	testDownloadDir  = "/downloads"
	testRelativePath = "relative/path"
	testTorrentName  = "test-torrent"
	testMagnetURI    = "magnet:?xt=urn:btih:abc123"
	testQueueTop     = "top"
	testActionAccept = "accept"
)

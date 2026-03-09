# mcp-transmission

MCP (Model Context Protocol) server for managing [Transmission](https://transmissionbt.com/) BitTorrent client.

Built on top of [go-transmission](https://github.com/lexfrei/go-transmission) library.

## Features

18 MCP tools covering the full Transmission RPC API:

**Torrent Management:**

- `transmission_torrent_list` — List all torrents with status, progress, speed
- `transmission_torrent_add` — Add torrent by magnet link, URL, or base64 .torrent
- `transmission_torrent_remove` — Remove torrents (optionally delete files)
- `transmission_torrent_start` — Start torrents (with queue bypass option)
- `transmission_torrent_stop` — Stop/pause torrents
- `transmission_torrent_verify` — Verify local data integrity
- `transmission_torrent_reannounce` — Force tracker announce
- `transmission_torrent_details` — Detailed info: files, trackers, peers
- `transmission_torrent_set` — Modify torrent properties (limits, labels, etc.)
- `transmission_torrent_move` — Move torrent data to new location

**Session Management:**

- `transmission_session_stats` — Session statistics and transfer totals
- `transmission_session_get` — Session configuration
- `transmission_session_set` — Modify session settings (speed limits, directories, etc.)

**System:**

- `transmission_free_space` — Check disk space
- `transmission_port_test` — Test peer port accessibility
- `transmission_blocklist_update` — Update IP blocklist

**Queue & Bandwidth:**

- `transmission_queue_move` — Move torrents in queue (top/up/down/bottom)
- `transmission_bandwidth_group_get` — Get bandwidth group configurations

## Configuration

| Variable | Description | Default |
| --- | --- | --- |
| `TRANSMISSION_URL` | Transmission RPC endpoint | `http://localhost:9091/transmission/rpc` |
| `TRANSMISSION_USERNAME` | HTTP Basic Auth username | (none) |
| `TRANSMISSION_PASSWORD` | HTTP Basic Auth password | (none) |
| `MCP_HTTP_PORT` | Optional HTTP/SSE transport port | (disabled) |

## Usage

### With Claude Code (stdio)

Add to your `.mcp.json`:

```json
{
  "mcp-transmission": {
    "command": "docker",
    "args": [
      "run", "--rm", "-i",
      "-e", "TRANSMISSION_URL",
      "-e", "TRANSMISSION_USERNAME",
      "-e", "TRANSMISSION_PASSWORD",
      "ghcr.io/lexfrei/mcp-transmission:latest"
    ],
    "env": {
      "TRANSMISSION_URL": "http://your-transmission-host:9091/transmission/rpc"
    }
  }
}
```

### Direct binary

```bash
TRANSMISSION_URL=http://localhost:9091/transmission/rpc ./mcp-transmission
```

### Container

> **Note:** CI/CD for publishing container images is not yet configured.
> Build the image locally using the instructions below.

```bash
docker build --file Containerfile --tag mcp-transmission .
docker run --rm -i \
  -e TRANSMISSION_URL=http://host.docker.internal:9091/transmission/rpc \
  mcp-transmission
```

## Building

```bash
go build ./cmd/mcp-transmission
```

### Container image

```bash
docker build --file Containerfile --tag mcp-transmission .
```

## License

BSD 3-Clause License. See [LICENSE](LICENSE).

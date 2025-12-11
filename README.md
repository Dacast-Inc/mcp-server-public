# Dacast MCP Server

This repository provides a **Model Context Protocol (MCP) server** that exposes the Dacast video platform API as a set of structured MCP tools. It lets MCP-compatible clients (such as Claude Desktop) manage Dacast resources — channels, playlists, images, simulcast destinations, and more — through natural language.

---

## Features

- **StdIO-based MCP server**
  - Speaks the [Model Context Protocol](https://modelcontextprotocol.io/) over standard input/output.
  - Designed to be launched by an MCP client process and kept running as a child process.

- **Dacast API integration**
  - Channel management: create, list, get, update.
  - Playlist management: create, list, get, update, set playlist content.
  - Image management: thumbnails and splash images.
  - Simulcast destination management: create, get, delete.

---

## Installation

### Prerequisites

- **Go**: Go 1.24+ (earlier versions may work but are not guaranteed).
- **Dacast account and API key**: required to perform authenticated operations.

### MCP Client Configuration

JSON based MCP client configuration might look like:

```json
{
  "mcpServers": {
    "dacast": {
      "command": "go",
      "args": ["run", "github.com/Dacast-Inc/mcp-server-public@latest"],
      "env": {
        "DACAST_API_KEY": "DACAST API KEY HERE"
      }
    }
  }
}
```

## Architecture Overview

### High-level design

The server is a single Go binary that:

1. Starts an MCP stdio server.
2. Registers a set of tools grouped by Dacast domain (channels, playlists, images, simulcast).
3. For each incoming MCP `call_tool` request:
   - Binds and validates the tool arguments.
   - Constructs an HTTP request to the relevant Dacast REST endpoint using the internal `ApiClient`.
   - Forwards the request to Dacast with the appropriate authentication headers.
   - Returns the JSON response (and optionally a transformed, structured form) back to the MCP client.

Key packages:

- `main.go` – entrypoint that initializes the MCP server and registers tools.
- `pkg/apiclient/` – minimal HTTP client wrapper around the Dacast REST API.
- `pkg/tools/` – tool implementations grouped by domain:
  - `channel/` – channel-related operations.
  - `playlist/` – playlist-related operations.
  - `images/` – thumbnail and splash image operations.
  - `simulcast/` – simulcast destination operations.
- `pkg/tools/toolscommon/` – shared types, handlers and utilities for building tools.

### Tools → Dacast API mapping

Each tool under `pkg/tools/` corresponds to one or more Dacast API endpoints.

- **Channels** (`pkg/tools/channel/`)
  - `create_channel.go` – create a new channel.
  - `get_channel.go` – get channel details.
  - `update_channel.go` – update an existing channel.
  - `list_channel.go` – list all or filtered channels.

- **Playlists** (`pkg/tools/playlist/`)
  - `create_playlist.go` – create a new playlist.
  - `get_playlist.go` – get playlist details.
  - `update_playlist.go` – update an existing playlist.
  - `list_playlist.go` – list playlists.
  - `set_playlist_content.go` – set or update playlist content.

- **Images** (`pkg/tools/images/`)
  - `thumbnail.go` – manage channel/asset thumbnails.
  - `splash.go` – manage splash images.

- **Simulcast** (`pkg/tools/simulcast/`)
  - `create_simulcast_destination.go` – create a simulcast destination.
  - `get_simulcast_destination.go` – retrieve a simulcast destination.
  - `delete_simulcast_destination.go` – delete a simulcast destination.
---

## License

This project is licensed under the terms described in the [`LICENSE`](./LICENSE) file in this repository. Please review that file for the full text.
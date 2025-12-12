package tools

import (
	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/channel"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/images"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/playlist"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/simulcast"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/vod"
	"github.com/mark3labs/mcp-go/server"
)

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	channel.Register(srv, client)
	simulcast.Register(srv, client)
	playlist.Register(srv, client)
	images.Register(srv, client)
	vod.Register(srv, client)
}

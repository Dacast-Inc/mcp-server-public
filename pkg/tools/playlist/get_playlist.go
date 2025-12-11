package playlist

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GetPlaylistRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Playlist ID to retrieve." jsonschema:"required"`
}

func (l GetPlaylistRequest) Transform() toolscommon.Transformable { return l }

func RegisterGetPlaylist(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("get_playlist_v1",
		mcp.WithDescription("Get playlist by ID. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[PlaylistResponse]()),
		mcp.WithInputSchema[GetPlaylistRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[GetPlaylistRequest, PlaylistResponse](
		client, http.MethodGet, "/v2/playlists/{id}", []string{"id"},
	))
}

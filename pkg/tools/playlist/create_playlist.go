package playlist

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CreatePlaylistRequest struct {
	Title           string  `json:"title" jsonschema_description:"Title of the playlist." jsonschema:"required"`
	Description     *string `json:"description,omitempty" jsonschema_description:"Description of the playlist."`
	Online          bool    `json:"online,omitempty" jsonschema_description:"Indicates if the playlist should be enabled (online) or disabled (offline) upon creation." jsonschema:"default=true"`
	GoogleAnalytics *string `json:"google_analytics,omitempty" jsonschema_description:"Google Analytics tracking code to associate with the playlist. If not set, tracking will be disabled."`
	Password        *string `json:"password,omitempty" jsonschema_description:"Password to protect the playlist. If not set, the playlist will not be password protected."`
}

func (l CreatePlaylistRequest) Transform() toolscommon.Transformable { return l }

func RegisterCreatePlaylist(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("create_playlist_v1",
		mcp.WithDescription("Create playlist. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[PlaylistResponse]()),
		mcp.WithInputSchema[CreatePlaylistRequest](),
		//mcp.WithOutputSchema[PlaylistResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[CreatePlaylistRequest, PlaylistResponse](
		client, http.MethodPost, "/v2/playlists", nil,
	))
}

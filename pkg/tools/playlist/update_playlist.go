package playlist

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type UpdatePlaylistRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Playlist ID to update." jsonschema:"required"`

	Title           *string `json:"title,omitempty" jsonschema_description:"Title of the playlist."`
	Description     *string `json:"description,omitempty" jsonschema_description:"Description of the playlist."`
	Online          bool    `json:"online,omitempty" jsonschema_description:"Indicates if the playlist should be enabled (online) or disabled (offline) upon creation." jsonschema:"default=true"`
	GoogleAnalytics *string `json:"google_analytics,omitempty" jsonschema_description:"Google Analytics tracking code to associate with the playlist. If not set, tracking will be disabled."`
	Password        *string `json:"password,omitempty" jsonschema_description:"Password to protect the playlist. If not set, the playlist will not be password protected."`
}

func (l UpdatePlaylistRequest) Transform() toolscommon.Transformable { return l }

func RegisterUpdatePlaylist(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("update_playlist_v1",
		mcp.WithDescription("Update playlist. You MUST obtain information about playlist and find out the values of the online field before calling this function using get_playlist or list_playlists if you do not intend to change them. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[PlaylistResponse]()),
		mcp.WithInputSchema[UpdatePlaylistRequest](),
		//mcp.WithOutputSchema[PlaylistResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[UpdatePlaylistRequest, PlaylistResponse](
		client, http.MethodPut, "/v2/playlists/{id}", []string{"id"},
	))
}

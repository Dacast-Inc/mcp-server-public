package playlist

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListPlaylistRequest struct {
	Page      *int    `json:"page,omitempty" url:"page,omitempty" jsonschema_description:"Page number for paginated results."`
	PerPage   *int    `json:"per_page,omitempty" url:"per_page,omitempty" jsonschema_description:"Number of items per page for paginated results." jsonschema:"minimum=1,maximum=100,default=50"`
	Title     *string `json:"title,omitempty" url:"title,omitempty" jsonschema_description:"Search for specific live streams using a substring or complete title"`
	StartDate *string `json:"start_date,omitempty" url:"start_date,omitempty" jsonschema_description:"Filter playlist created after this date. YYYY-MM-DD format."`
	EndDate   *string `json:"end_date,omitempty" url:"end_date,omitempty" jsonschema_description:"Filter playlist created before this date. YYYY-MM-DD format."`
}

func (l ListPlaylistRequest) Transform() toolscommon.Transformable { return l }

func RegisterListPlaylist(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("list_playlists_v1",
		mcp.WithDescription("Returns a paginated list of all playlists. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[toolscommon.PaginatedData[PlaylistResponse]]()),
		mcp.WithInputSchema[ListPlaylistRequest](),
		//mcp.WithOutputSchema[toolscommon.PaginatedData[ChannelResponse]](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[ListPlaylistRequest, toolscommon.PaginatedData[PlaylistResponse]](
		client, http.MethodGet, "/v2/playlists", nil,
	))
}

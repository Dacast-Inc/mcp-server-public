package channel

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListChannelRequest struct {
	Page                *int    `json:"page,omitempty" url:"page,omitempty" jsonschema_description:"Page number for paginated results."`
	PerPage             *int    `json:"per_page,omitempty" url:"per_page,omitempty" jsonschema_description:"Number of items per page for paginated results." jsonschema:"minimum=1,maximum=100,default=50"`
	Title               *string `json:"title,omitempty" url:"title,omitempty" jsonschema_description:"Search for specific live streams using a substring or complete title"`
	StartDate           *string `json:"start_date,omitempty" url:"start_date,omitempty" jsonschema_description:"Filter channels created after this date. YYYY-MM-DD format."`
	EndDate             *string `json:"end_date,omitempty" url:"end_date,omitempty" jsonschema_description:"Filter channels created before this date. YYYY-MM-DD format."`
	RequestOnlineStatus *bool   `json:"request_online_status,omitempty" url:"request_online_status,omitempty" jsonschema_description:"If true, the response will include the online status of each channel. This may increase response time. Set this flag only if it is really necessary for processing the user request." jsonschema:"default=false"`
}

func (l ListChannelRequest) Transform() toolscommon.Transformable { return l }

func RegisterListChannel(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("list_channels_v1",
		mcp.WithDescription("Returns a paginated list of all channels (streams). \n\nResponse schema: "+toolscommon.TypeToJsonSchema[toolscommon.PaginatedData[ChannelResponse]]()),
		mcp.WithInputSchema[ListChannelRequest](),
		//mcp.WithOutputSchema[toolscommon.PaginatedData[ChannelResponse]](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[ListChannelRequest, toolscommon.PaginatedData[ChannelResponse]](
		client, http.MethodGet, "/v2/channel", nil,
	))
}

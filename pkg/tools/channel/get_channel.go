package channel

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GetChannelRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Channel (stream) ID to retrieve." jsonschema:"required"`
}

func (l GetChannelRequest) Transform() toolscommon.Transformable { return l }

func RegisterGetChannel(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("get_channel_v1",
		mcp.WithDescription("Get channel (stream) by ID. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[ChannelResponse]()),
		mcp.WithInputSchema[GetChannelRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[GetChannelRequest, ChannelResponse](
		client, http.MethodGet, "/v2/channel/{id}", []string{"id"},
	))
}

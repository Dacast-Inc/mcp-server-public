package simulcast

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CreateSimulcastDestinationRequest struct {
	ID        string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Channel (stream) ID to add the simulcast destination to." jsonschema:"required"`
	RtmpUrl   string `json:"rtmp_url,omitempty" jsonschema_description:"RTMP url for the simulcast destination (stream)." jsonschema:"required"`
	StreamKey string `json:"stream_key,omitempty" jsonschema_description:"Stream key for the simulcast destination (stream)." jsonschema:"required"`
}

func (l CreateSimulcastDestinationRequest) Transform() toolscommon.Transformable { return l }

func RegisterCreateSimulcastDestination(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("create_simulcast_destination_v1",
		mcp.WithDescription("Add a simulcast destination to a channel (stream) \n\nResponse schema: "+toolscommon.TypeToJsonSchema[EmptyObjectResponse]()),
		mcp.WithInputSchema[CreateSimulcastDestinationRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[CreateSimulcastDestinationRequest, EmptyObjectResponse](
		client, http.MethodPost, "/v2/channel/{id}/simulcast-destinations", []string{"id"},
	))
}

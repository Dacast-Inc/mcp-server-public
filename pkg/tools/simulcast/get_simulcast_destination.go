package simulcast

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GetSimulcastDestinationsRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Channel (stream) ID to retrieve the simulcast destinations from." jsonschema:"required"`
}

func (l GetSimulcastDestinationsRequest) Transform() toolscommon.Transformable { return l }

func RegisterGetSimulcastDestinations(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("get_simulcast_destinations_v1",
		mcp.WithDescription("List the simulcast destinations of a channel (stream). \n\nResponse schema: "+toolscommon.TypeToJsonSchema[SimulcastResponse]()),
		mcp.WithInputSchema[GetSimulcastDestinationsRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[GetSimulcastDestinationsRequest, SimulcastResponse](
		client, http.MethodGet, "/v2/channel/{id}/simulcast-destinations", []string{"id"},
	))
}

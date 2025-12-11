package simulcast

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DeleteSimulcastDestinationRequest struct {
	ID            string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Channel (stream) ID to remove the simulcast destination from." jsonschema:"required"`
	DestinationID string `json:"destination_id,omitempty" jsonschema_description:"ID of the simulcast destination." jsonschema:"required"`
}

func (l DeleteSimulcastDestinationRequest) Transform() toolscommon.Transformable { return l }

func RegisterDeleteSimulcastDestination(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("delete_simulcast_destination_v1",
		mcp.WithDescription("Remove a simulcast destination from the channel (stream) \n\nResponse schema: "+toolscommon.TypeToJsonSchema[EmptyObjectResponse]()),
		mcp.WithInputSchema[DeleteSimulcastDestinationRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[DeleteSimulcastDestinationRequest, EmptyObjectResponse](
		client, http.MethodDelete, "/v2/channel/{id}/simulcast-destinations/{destination_id}", []string{"id", "destination_id"},
	))
}

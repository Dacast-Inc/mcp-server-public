package simulcast

import (
	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/server"
)

type EmptyObjectResponse struct {
}

func (c EmptyObjectResponse) Transform() toolscommon.Transformable {
	return c
}

type SimulcastResponse struct {
	Data []SimulcastDestination `json:"data" jsonschema_description:"List of simulcast destinations associated with the channel (stream)."`
}

type SimulcastDestination struct {
	ID        string `json:"id" jsonschema_description:"Unique identifier of the simulcast destination."`
	RtmpUrl   string `json:"rtmp_url" jsonschema_description:"RTMP url for the simulcast destination (stream)."`
	StreamKey string `json:"stream_key" jsonschema_description:"Stream key for the simulcast destination (stream)."`
}

func (c SimulcastResponse) Transform() toolscommon.Transformable {
	return c
}

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	RegisterGetSimulcastDestinations(srv, client)
	RegisterCreateSimulcastDestination(srv, client)
	RegisterDeleteSimulcastDestination(srv, client)
}

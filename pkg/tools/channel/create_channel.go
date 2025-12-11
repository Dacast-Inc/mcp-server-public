package channel

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CreateChannelRequest struct {
	Title                string  `json:"title" jsonschema_description:"Title of the channel (stream)." jsonschema:"required"`
	Description          *string `json:"description,omitempty" jsonschema_description:"Description of the channel (stream)."`
	Region               *string `json:"region,omitempty" jsonschema_description:"Region where the channel (stream) will be hosted. If you have information about the user's location, you should select the closest region. If not, you should ask the user about their region." jsonschema:"enum=north_america,enum=europe,enum=asia_pacific,default=north_america"`
	Online               *bool   `json:"online,omitempty" jsonschema_description:"Indicates if the channel (stream) should be enabled (online) or disabled (offline) upon creation." jsonschema:"default=true"`
	LiveRecordingEnabled *bool   `json:"live_recording_enabled,omitempty" jsonschema_description:"Indicates if live recording enabled for the channel (stream) upon creation." jsonschema:"default=false"`
	LiveDVREnabled       *bool   `json:"live_dvr_enabled,omitempty" jsonschema_description:"Indicates if live DVR is enabled for the channel (stream) upon creation. Cannot be enabled without recording" jsonschema:"default=false"`
	GoogleAnalyticsCode  *string `json:"google_analytics_code,omitempty" jsonschema_description:"Google Analytics tracking code to associate with the channel (stream). If not set, tracking will be disabled."`
	Password             *string `json:"password,omitempty" jsonschema_description:"Password to protect the channel (stream). If not set, the channel will not be password protected."`
	ChannelType          *string `json:"channel_type,omitempty" jsonschema_description:"Type of the channel (stream). Transmux - just passthrough video/audio data. fhd-transcode - transcode video data and create additional ABR variants, If ABR is enabled and the stream has fewer than 5 viewers, an additional fee will be charged. HLS - Instead of rtmp, channel provides a direct endpoint for HLS publishing." jsonschema:"enum=transmux,enum=fhd-transcode,enum=hls,default=transmux"`
}

func (l CreateChannelRequest) Transform() toolscommon.Transformable { return l }

func RegisterCreateChannel(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("create_channel_v1",
		mcp.WithDescription("Create channel (stream). If you create an ABR channel, you are required to notify users of an additional fee for channels with fewer than 5 viewers. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[ChannelResponse]()),
		mcp.WithInputSchema[CreateChannelRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[CreateChannelRequest, ChannelResponse](
		client, http.MethodPost, "/v2/channel", nil,
	))
}

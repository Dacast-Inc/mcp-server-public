package channel

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type UpdateChannelRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Channel (stream) ID to update." jsonschema:"required"`

	Title                *string `json:"title,omitempty" jsonschema_description:"Title of the channel (stream)."`
	Description          *string `json:"description,omitempty" jsonschema_description:"Description of the channel (stream)."`
	Online               *bool   `json:"online,omitempty" jsonschema_description:"Indicates if the channel (stream) should be enabled (online) or disabled (offline)" jsonschema:"default=true"`
	LiveRecordingEnabled *bool   `json:"live_recording_enabled,omitempty" jsonschema_description:"Indicates if live recording enabled for the channel (stream)" jsonschema:"default=false"`
	LiveDVREnabled       *bool   `json:"live_dvr_enabled,omitempty" jsonschema_description:"Indicates if live DVR is enabled for the channel (stream). Cannot be enabled without recording" jsonschema:"default=false"`
	GoogleAnalyticsCode  *string `json:"google_analytics_code,omitempty" jsonschema_description:"Google Analytics tracking code to associate with the channel (stream). If not set, tracking will be disabled."`
	Password             *string `json:"password,omitempty" jsonschema_description:"Password to protect the channel (stream). If not set, the channel will not be password protected."`
}

func (l UpdateChannelRequest) Transform() toolscommon.Transformable { return l }

func RegisterUpdateChannel(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("update_channel_v1",
		mcp.WithDescription("Update channel (stream). You MUST obtain information about the stream and find out the values of the online, live_recording_enabled, and live_dvr_enabled fields before calling this function using get channel or list_channels if you do not intend to change them. Some fields, such as recording and dvr, are not updated instantly, so the resulting object may contain incorrect information about this state immediately after the update. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[ChannelResponse]()),
		mcp.WithInputSchema[UpdateChannelRequest](),
		//mcp.WithOutputSchema[ChannelResponse](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[UpdateChannelRequest, ChannelResponse](
		client, http.MethodPut, "/v2/channel/{id}", []string{"id"},
	))
}

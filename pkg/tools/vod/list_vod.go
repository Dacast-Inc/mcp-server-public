package vod

import (
	"net/http"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListVodRequest struct {
	Page        *int    `json:"page,omitempty" url:"page,omitempty" jsonschema_description:"Page number for paginated results."`
	PerPage     *int    `json:"per_page,omitempty" url:"per_page,omitempty" jsonschema_description:"Number of items per page for paginated results." jsonschema:"minimum=1,maximum=100,default=50"`
	Title       *string `json:"title,omitempty" url:"title,omitempty" jsonschema_description:"Search for specific VODs using a substring or complete title"`
	StreamID    *string `json:"stream_id,omitempty" url:"stream_id,omitempty" jsonschema_description:"Filter VODs by parent live channel ID (stream ID). Only returns VODs that were created from recordings of the specified live channel."`
	DurationMin *int    `json:"duration_min,omitempty" url:"duration_min,omitempty" jsonschema_description:"Filter VODs by minimum duration in seconds. Returns VODs with duration greater than or equal to this value. Example: 600 for videos longer than 10 minutes."`
	DurationMax *int    `json:"duration_max,omitempty" url:"duration_max,omitempty" jsonschema_description:"Filter VODs by maximum duration in seconds. Returns VODs with duration less than or equal to this value."`
	StartDate   *string `json:"start_date,omitempty" url:"start_date,omitempty" jsonschema_description:"Filter VODs created after this date. YYYY-MM-DD format."`
	EndDate     *string `json:"end_date,omitempty" url:"end_date,omitempty" jsonschema_description:"Filter VODs created before this date. YYYY-MM-DD format."`
}

func (l ListVodRequest) Transform() toolscommon.Transformable { return l }

func RegisterListVod(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("list_vods_v1",
		mcp.WithDescription("Returns a paginated list of all VODs. Can filter by stream_id, duration (duration_min/duration_max), or date range (start_date/end_date) to get VODs matching specific criteria. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[toolscommon.PaginatedData[VodResponse]]()),
		mcp.WithInputSchema[ListVodRequest](),
		//mcp.WithOutputSchema[toolscommon.PaginatedData[VodResponse]](),
	)

	srv.AddTool(tool, toolscommon.DefaultDacastProxyHandler[ListVodRequest, toolscommon.PaginatedData[VodResponse]](
		client, http.MethodGet, "/v2/vod", nil,
	))
}

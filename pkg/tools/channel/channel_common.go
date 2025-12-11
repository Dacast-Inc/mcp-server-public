package channel

import (
	"strings"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/server"
)

type ChannelResponse struct {
	ID          string  `json:"id" jsonschema_description:"Unique identifier of the channel."`
	Title       string  `json:"title" jsonschema_description:"Title of the channel."`
	Description *string `json:"description,omitempty" jsonschema_description:"Description of the channel."`

	AssetID    string `json:"asset_id" jsonschema_description:"Also known as contentId. Used by playback api to generate hls playback urls."`
	Autoplay   bool   `json:"autoplay" jsonschema_description:"Indicates if the channel should autoplay the video when loaded."`
	CompanyUrl string `json:"company_url" jsonschema_description:"URL of the company owning the channel."`

	Config struct {
		PublishingPointPrimary *string `json:"publishing_point_primary,omitempty" jsonschema_description:"Primary publish point URL. Has rtmp:// prefix for RTMP streams, and https:// prefix direct HLS streams. If not set, then its VOD To Live channel."`
		PublishingPointBackup  *string `json:"publishing_point_backup,omitempty" jsonschema_description:"Backup publish point URL. May be null if not configured."`
		StreamName             *string `json:"stream_name,omitempty" jsonschema_description:"Its actually stream key. Used for RTMP stream publishing."`
	} `json:"config" jsonschema_description:"Encoder configuration settings for the channel."`

	CountdownDate      *string `json:"countdown_date,omitempty" jsonschema_description:"If set, indicates the date and time when the countdown ends and the channel goes live. ISO 8601, but timezone is always Z. Always UTC time."`
	CreationDate       string  `json:"creation_date" jsonschema_description:"Channel creation date in ISO 8601 format. UTC timezone."`
	EnableCoupon       bool    `json:"enable_coupon" jsonschema_description:"Indicates if coupon codes are enabled for the channel. Paywall channels only."`
	EnablePayperview   bool    `json:"enable_payperview" jsonschema_description:"Indicates if pay-per-view is enabled for the channel. Paywall channels only."`
	EnableSubscription bool    `json:"enable_subscription" jsonschema_description:"Indicates if subscriptions is enabled for the channel. Paywall channels only."`
	Hls                string  `json:"hls" jsonschema_description:"HLS playback URL for the channel. This is not the preferred method for embedding content. Only use this field if the user needs to embed content in their own player or use it in their own software. The playback URL returned by this API is unsecure: it never expires and can be copy pasted by viewers to redistribute the stream. If user want a secure playback URL user can use the Playback API https://docs.dacast.com/reference/generate-hls-playback-url"`

	LiveRecordingEnabled bool    `json:"live_recording_enabled" jsonschema_description:"Indicates if live recordings are enabled for the channel."`
	LiveDVREnabled       bool    `json:"live_dvr_enabled" jsonschema_description:"Indicates if live DVR is enabled for the channel. Cannot be enabled if live recording is disabled."`
	Online               *bool   `json:"online" jsonschema_description:"Indicates if the channel is currently enabled or disabled. Disabled channels cannot be viewed by end users. THIS FIELD DOES NOT SHOW WHETHER THE CHANNEL IS ONLINE, ONLY ON/OFF. Only use this field if the user asks whether the channel is enabled, but do not use this field in the context of whether it is online. Do not explain the specifics of this field unless you are explicitly asked to do so."`
	OnlineSince          *string `json:"online_since,omitempty" jsonschema_description:"If set, then channel is currently online (streaming), indicates the date and time when the channel went online. Unix timestamp format. UTC timezone."`
	OnlineSinceDate      *string `json:"online_since_date,omitempty" jsonschema_description:"If the channel is currently online (streaming), indicates the date when the channel went online. ISO 8601 format. UTC timezone."`
	Password             *string `json:"password,omitempty" jsonschema_description:"If the channel is password protected, indicates the password required to view the channel."`
	SaveDate             string  `json:"save_date" jsonschema_description:"Date when the channel was last saved. ISO 8601 format. UTC timezone."`
	ShareCode            *struct {
		Twitter string `json:"twitter"` // just to extract the link
	} `json:"share_code,omitempty"`

	ShareLink string `json:"share_link" jsonschema_description:"Public URL of the channel. Can be used to share the channel with viewers or to embed the channel in a webpage (iframe). This is the preferred method for embedding content."`
}

func (c ChannelResponse) Transform() toolscommon.Transformable {
	if c.ShareCode != nil {
		c.ShareLink = c.ShareCode.Twitter
		c.ShareCode = nil
	} else if c.AssetID != "" {
		parts := strings.Split(c.AssetID, "-live-")
		if len(parts) >= 3 {
			c.ShareLink = "https://iframe.dacast.com/live/" + parts[0] + "/" + parts[2]
		}
	}

	return c
}

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	RegisterListChannel(srv, client)
	RegisterGetChannel(srv, client)
	RegisterCreateChannel(srv, client)
	RegisterUpdateChannel(srv, client)
}

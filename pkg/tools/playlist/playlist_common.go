package playlist

import (
	"strings"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/server"
)

type PlaylistResponse struct {
	ID          string  `json:"id" jsonschema_description:"Unique identifier of the playlist."`
	Title       string  `json:"title" jsonschema_description:"Title of the playlist."`
	Description *string `json:"description,omitempty" jsonschema_description:"Description of the playlist."`

	AssetID    string `json:"asset_id" jsonschema_description:"Also known as contentId. Used by playback api to generate hls playback urls."`
	Autoplay   bool   `json:"autoplay" jsonschema_description:"Indicates if the playlist should autoplay the video when loaded."`
	CompanyUrl string `json:"company_url" jsonschema_description:"URL of the company owning the playlist."`

	Content *struct {
		List []PlaylistContent `json:"list" jsonschema_description:"List of content items in the playlist."`
	}

	CountdownDate *string `json:"countdown_date,omitempty" jsonschema_description:"If set, indicates the date and time when the countdown ends and the channel goes live. ISO 8601, but timezone is always Z. Always UTC time."`
	CreationDate  string  `json:"creation_date" jsonschema_description:"playlist creation date in ISO 8601 format. UTC timezone."`

	EnablePayperview   bool `json:"enable_payperview" jsonschema_description:"Indicates if pay-per-view is enabled for the playlist. Paywall playlists only."`
	EnableSubscription bool `json:"enable_subscription" jsonschema_description:"Indicates if subscriptions is enabled for the playlist. Paywall playlists only."`

	Online   *bool   `json:"online" jsonschema_description:"Indicates if the playlist is currently enabled or disabled. Disabled playlist cannot be viewed by end users. THIS FIELD DOES NOT SHOW WHETHER THE CHANNEL IS ONLINE, ONLY ON/OFF. Only use this field if the user asks whether the channel is enabled, but do not use this field in the context of whether it is online. Do not explain the specifics of this field unless you are explicitly asked to do so."`
	Password *string `json:"password,omitempty" jsonschema_description:"If the playlist is password protected, indicates the password required to view the channel."`
	SaveDate string  `json:"save_date" jsonschema_description:"Date when the playlist was last saved. ISO 8601 format. UTC timezone."`

	ShareCode *struct {
		Twitter string `json:"twitter"` // just to extract the link
	} `json:"share_code,omitempty"`

	ShareLink string `json:"share_link" jsonschema_description:"Public URL of the playlist. Can be used to share the playlist with viewers or to embed the playlist in a webpage (iframe). This is the preferred method for embedding content."`
}

func (c PlaylistResponse) Transform() toolscommon.Transformable {
	if c.ShareCode != nil {
		c.ShareLink = c.ShareCode.Twitter
		c.ShareCode = nil
	} else if c.AssetID != "" {
		parts := strings.Split(c.AssetID, "-playlist-")
		if len(parts) >= 2 {
			c.ShareLink = "https://iframe.dacast.com/playlist/" + parts[0] + "/" + parts[1]
		}
	}

	return c
}

type PlaylistContent struct {
	ID       string `json:"id" jsonschema_description:"Unique identifier of the playlist content item (VOD id or channel id)."`
	Position string `json:"position" jsonschema_description:"Position of the playlist content."`
	Title    string `json:"title" jsonschema_description:"Title of the playlist content item."`
	Type     string `json:"type" jsonschema_description:"Type of the playlist content item. Can be 'vod' or 'channel'."`
}

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	RegisterListPlaylist(srv, client)
	RegisterGetPlaylist(srv, client)
	RegisterCreatePlaylist(srv, client)
	RegisterUpdatePlaylist(srv, client)
	RegisterAddPlaylistContent(srv, client)
	RegisterSetPlaylistContentPosition(srv, client)
}

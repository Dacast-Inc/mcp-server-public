package vod

import (
	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/server"
)

type VodResponse struct {
	Abitrate    *int    `json:"abitrate,omitempty" mapstructure:"abitrate" jsonschema_description:"Audio bitrate in bits per second."`
	Acodec      *string `json:"acodec,omitempty" mapstructure:"acodec" jsonschema_description:"Audio codec used (e.g., aac)."`
	Vbitrate    *int    `json:"vbitrate,omitempty" mapstructure:"vbitrate" jsonschema_description:"Video bitrate in bits per second."`
	Vcodec      *string `json:"vcodec,omitempty" mapstructure:"vcodec" jsonschema_description:"Video codec used (e.g., h264)."`
	VideoWidth  *int    `json:"video_width,omitempty" mapstructure:"video_width" jsonschema_description:"Width of the video in pixels."`
	VideoHeight *int    `json:"video_height,omitempty" mapstructure:"video_height" jsonschema_description:"Height of the video in pixels."`
	Container   *string `json:"container,omitempty" mapstructure:"container" jsonschema_description:"Container format (e.g., mpeg-4)."`

	ID          string  `json:"id" mapstructure:"id" jsonschema_description:"Unique identifier of the VOD."`
	OriginalID  *string `json:"original_id,omitempty" mapstructure:"original_id" jsonschema_description:"Original ID of the VOD."`
	Title       string  `json:"title" mapstructure:"title" jsonschema_description:"Title of the VOD."`
	Description *string `json:"description,omitempty" mapstructure:"description" jsonschema_description:"Description of the VOD."`

	AssetID   string     `json:"asset_id" mapstructure:"asset_id" jsonschema_description:"Also known as contentId. Used by playback api to generate hls playback urls."`
	Autoplay  bool       `json:"autoplay" mapstructure:"autoplay" jsonschema_description:"Indicates if the VOD should autoplay the video when loaded."`
	Hls       *string    `json:"hls,omitempty" mapstructure:"hls" jsonschema_description:"HLS playback URL for the VOD. This is not the preferred method for embedding content. Only use this field if the user needs to embed content in their own player or use it in their own software. The playback URL returned by this API is unsecure: it never expires and can be copy pasted by viewers to redistribute the stream. If user want a secure playback URL user can use the Playback API https://docs.dacast.com/reference/generate-hls-playback-url"`
	Hds       *string    `json:"hds,omitempty" mapstructure:"hds" jsonschema_description:"HDS playback URL for the VOD."`
	ShareCode *ShareCode `json:"share_code,omitempty" mapstructure:"share_code" jsonschema_description:"Share links for the VOD (facebook, twitter, gplus)."`

	Online       bool   `json:"online" mapstructure:"online" jsonschema_description:"Indicates if the VOD is currently enabled or disabled. Disabled VODs cannot be viewed by end users."`
	IsSecured    *bool  `json:"is_secured,omitempty" mapstructure:"is_secured" jsonschema_description:"Indicates if the VOD is secured."`
	SaveDate     string `json:"save_date" mapstructure:"save_date" jsonschema_description:"Date when the VOD was last saved. ISO 8601 format. UTC timezone."`
	CreationDate string `json:"creation_date" mapstructure:"creation_date" jsonschema_description:"VOD creation date in ISO 8601 format. UTC timezone."`

	Password        *string `json:"password,omitempty" mapstructure:"password" jsonschema_description:"If the VOD is password protected, indicates the password required to view the VOD."`
	NoframeSecurity *int    `json:"noframe_security,omitempty" mapstructure:"noframe_security" jsonschema_description:"No frame security setting."`

	PaywallEnabled     bool `json:"paywall-enabled" mapstructure:"paywall-enabled" jsonschema_description:"Indicates if paywall is enabled for the VOD."`
	EnableCoupon       bool `json:"enable_coupon" mapstructure:"enable_coupon" jsonschema_description:"Indicates if coupon codes are enabled for the VOD. Paywall VODs only."`
	EnablePayperview   bool `json:"enable_payperview" mapstructure:"enable_payperview" jsonschema_description:"Indicates if pay-per-view is enabled for the VOD. Paywall VODs only."`
	EnableSubscription bool `json:"enable_subscription" mapstructure:"enable_subscription" jsonschema_description:"Indicates if subscriptions is enabled for the VOD. Paywall VODs only."`

	Filename   *string `json:"filename,omitempty" mapstructure:"filename" jsonschema_description:"Name of the uploaded file."`
	Filesize   *string `json:"filesize,omitempty" mapstructure:"filesize" jsonschema_description:"Size of the file in bytes (as string)."`
	DiskUsage  *int64  `json:"disk_usage,omitempty" mapstructure:"disk_usage" jsonschema_description:"Disk usage in bytes."`
	Duration   *string `json:"duration,omitempty" mapstructure:"duration" jsonschema_description:"Duration of the video in HH:MM:SS format."`
	Streamable *int    `json:"streamable,omitempty" mapstructure:"streamable" jsonschema_description:"Streamable flag."`

	Renditions *[]Rendition `json:"renditions,omitempty" mapstructure:"renditions" jsonschema_description:"List of available renditions (quality variants) for the VOD."`

	CategoryID         *int      `json:"category_id,omitempty" mapstructure:"category_id" jsonschema_description:"Category ID for the VOD."`
	GroupID            *int      `json:"group_id,omitempty" mapstructure:"group_id" jsonschema_description:"Group ID for the VOD."`
	Folders            *[]Folder `json:"folders,omitempty" mapstructure:"folders" jsonschema_description:"Folders the VOD belongs to."`
	AssociatedPackages *string   `json:"associated_packages,omitempty" mapstructure:"associated_packages" jsonschema_description:"Associated packages."`

	CountriesID *string `json:"countries_id,omitempty" mapstructure:"countries_id" jsonschema_description:"Countries restriction ID."`
	ReferersID  *string `json:"referers_id,omitempty" mapstructure:"referers_id" jsonschema_description:"Referers restriction ID."`

	PlayerHeight *int `json:"player_height,omitempty" mapstructure:"player_height" jsonschema_description:"Player height in pixels."`
	PlayerWidth  *int `json:"player_width,omitempty" mapstructure:"player_width" jsonschema_description:"Player width in pixels."`

	PublishOnDacast   *bool   `json:"publish_on_dacast,omitempty" mapstructure:"publish_on_dacast" jsonschema_description:"Indicates if the VOD is published on Dacast."`
	ExternalVideoPage *string `json:"external_video_page,omitempty" mapstructure:"external_video_page" jsonschema_description:"External video page URL."`

	SplashscreenID *int      `json:"splashscreen_id,omitempty" mapstructure:"splashscreen_id" jsonschema_description:"Splashscreen asset ID."`
	ThumbnailID    *int      `json:"thumbnail_id,omitempty" mapstructure:"thumbnail_id" jsonschema_description:"Thumbnail asset ID."`
	Pictures       *Pictures `json:"pictures,omitempty" mapstructure:"pictures" jsonschema_description:"Picture URLs (splashscreen, thumbnail)."`

	Subtitles *[]Subtitle `json:"subtitles,omitempty" mapstructure:"subtitles" jsonschema_description:"Subtitles associated with the VOD."`

	TemplateID *int `json:"template_id,omitempty" mapstructure:"template_id" jsonschema_description:"Template ID for the VOD."`

	GoogleAnalytics *int `json:"google_analytics,omitempty" mapstructure:"google_analytics" jsonschema_description:"Google Analytics setting."`

	CustomData *string `json:"custom_data,omitempty" mapstructure:"custom_data" jsonschema_description:"Custom data associated with the VOD."`

	StreamID *string `json:"stream_id,omitempty" mapstructure:"stream_id" jsonschema_description:"Parent live channel ID if this VOD was created from a live stream recording."`

	Ads *string `json:"ads,omitempty" mapstructure:"ads" jsonschema_description:"Ads ID for the VOD."`
}

type ShareCode struct {
	Facebook *string `json:"facebook,omitempty" mapstructure:"facebook" jsonschema_description:"Facebook share link (URL) for the VOD."`
	Twitter  *string `json:"twitter,omitempty" mapstructure:"twitter" jsonschema_description:"Twitter (X.com) share link (URL) for the VOD."`
	Gplus    *string `json:"gplus,omitempty" mapstructure:"gplus" jsonschema_description:"Google Plus share link (URL) for the VOD."`
}

type Rendition struct {
	ID          string  `json:"id" mapstructure:"id" jsonschema_description:"Rendition ID."`
	Width       *int    `json:"width,omitempty" mapstructure:"width" jsonschema_description:"Width of the rendition in pixels."`
	Height      *int    `json:"height,omitempty" mapstructure:"height" jsonschema_description:"Height of the rendition in pixels."`
	Bitrate     *int    `json:"bitrate,omitempty" mapstructure:"bitrate" jsonschema_description:"Bitrate of the rendition in bits per second."`
	Framerate   *string `json:"framerate,omitempty" mapstructure:"framerate" jsonschema_description:"Framerate of the rendition (e.g., 60)."`
	SizeInBytes *int64  `json:"size_in_bytes,omitempty" mapstructure:"size_in_bytes" jsonschema_description:"Size of the rendition file in bytes."`
}

type Folder struct {
	Id       *string `json:"id,omitempty" mapstructure:"id" jsonschema_description:"Folder ID."`
	Path     *string `json:"path,omitempty" mapstructure:"path" jsonschema_description:"Path of the folder."`
	ParentId *string `json:"parent_id,omitempty" mapstructure:"parent-id" jsonschema_description:"Parent folder ID."`
	Name     *string `json:"name,omitempty" mapstructure:"name" jsonschema_description:"Name of the folder."`
}

type Pictures struct {
	Splashscreen *[]string `json:"splashscreen,omitempty" mapstructure:"splashscreen" jsonschema_description:"Splashscreen image URLs."`
	Thumbnail    *[]string `json:"thumbnail,omitempty" mapstructure:"thumbnail" jsonschema_description:"Thumbnail image URLs."`
}

type Subtitle struct {
	ID   *string `json:"id,omitempty" mapstructure:"id" jsonschema_description:"Subtitle ID."`
	URL  *string `json:"url,omitempty" mapstructure:"url" jsonschema_description:"Subtitle file URL."`
	Lang *string `json:"lang,omitempty" mapstructure:"lang" jsonschema_description:"Language code of the subtitle."`
}

func (c VodResponse) Transform() toolscommon.Transformable {
	return c
}

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	RegisterListVod(srv, client)
}

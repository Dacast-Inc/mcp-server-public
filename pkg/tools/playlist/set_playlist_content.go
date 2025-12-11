package playlist

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/gotidy/ptr"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type AddPlaylistContentRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Playlist ID to update." jsonschema:"required"`

	ContentID string `json:"content_id,omitempty" url:"content_id,omitempty" jsonschema_description:"Content ID to add to the playlist." jsonschema:"required"`
	Type      string `json:"type,omitempty" url:"type,omitempty" jsonschema_description:"Type of content to add to the playlist.." jsonschema:"required,enum=channel"`
	Position  int    `json:"position,omitempty" url:"position,omitempty" jsonschema_description:"Position in playlist. Positions after the specified one will be automatically shifted by +1" jsonschema:"default=1"`
}

type SetPlaylistContentPositionRequest struct {
	ID string `json:"id,omitempty" url:"id,omitempty" jsonschema_description:"Playlist ID to update." jsonschema:"required"`

	ContentID string `json:"content_id,omitempty" url:"content_id,omitempty" jsonschema_description:"Content ID to update position." jsonschema:"required"`
	Position  int    `json:"position,omitempty" url:"position,omitempty" jsonschema_description:"Position in playlist. Positions after the specified one will be automatically shifted by +1" jsonschema:"required"`
}

type ContentItem struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

type SetPlaylistContentRequest struct {
	Content   *string `json:"content,omitempty"`
	Positions *string `json:"positions,omitempty"`
}

func (l AddPlaylistContentRequest) Transform() toolscommon.Transformable         { return l }
func (l SetPlaylistContentPositionRequest) Transform() toolscommon.Transformable { return l }

func RegisterAddPlaylistContent(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("add_content_to_playlist_v1",
		mcp.WithDescription("Add content to playlist. Content may be added to playlist with some delay. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[PlaylistResponse]()),
		mcp.WithInputSchema[AddPlaylistContentRequest](),
		//mcp.WithOutputSchema[PlaylistResponse](),
	)

	srv.AddTool(tool, toolscommon.DacastHandler[AddPlaylistContentRequest, PlaylistResponse](func(logger *slog.Logger, args AddPlaylistContentRequest, headers map[string]string) ([]byte, error) {
		_, body, err := client.DoRequest(http.MethodGet, "/v2/playlists/"+args.ID+"", url.Values{}, nil, headers)
		if err != nil {
			return nil, err
		}

		var playlist PlaylistResponse
		if err = json.Unmarshal(body, &playlist); err != nil {
			return nil, err
		}

		content := appendPlaylistContent(&playlist, ContentItem{
			ID:       args.ContentID,
			Type:     args.Type,
			Position: args.Position,
		})

		contentBytes, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}

		req := SetPlaylistContentRequest{
			Content: ptr.Of(string(contentBytes)),
		}

		reqString, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		_, respBody, err := client.DoRequest(http.MethodPut, "/v2/playlists/"+args.ID+"/content", url.Values{}, ptr.Of(string(reqString)), headers)

		return respBody, err
	}))
}

func RegisterSetPlaylistContentPosition(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("set_playlist_content_position_v1",
		mcp.WithDescription("Set playlist content position. \n\nResponse schema: "+toolscommon.TypeToJsonSchema[PlaylistResponse]()),
		mcp.WithInputSchema[SetPlaylistContentPositionRequest](),
		//mcp.WithOutputSchema[PlaylistResponse](),
	)

	srv.AddTool(tool, toolscommon.DacastHandler[SetPlaylistContentPositionRequest, PlaylistResponse](func(logger *slog.Logger, args SetPlaylistContentPositionRequest, headers map[string]string) ([]byte, error) {
		_, body, err := client.DoRequest(http.MethodGet, "/v2/playlists/"+args.ID+"", url.Values{}, nil, headers)
		if err != nil {
			return nil, err
		}

		var playlist PlaylistResponse
		if err = json.Unmarshal(body, &playlist); err != nil {
			return nil, err
		}

		var contentItem *ContentItem

		for id, item := range playlist.Content.List {
			if item.ID == args.ContentID {
				contentItem = &ContentItem{
					ID:       item.ID,
					Type:     item.Type,
					Position: args.Position,
				}

				playlist.Content.List = append(playlist.Content.List[:id], playlist.Content.List[id+1:]...)

				break
			}
		}

		if contentItem == nil {
			return nil, fmt.Errorf("Invalid content ID")
		}

		content := appendPlaylistContent(&playlist, *contentItem)

		contentBytes, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}

		req := SetPlaylistContentRequest{
			Content: ptr.Of(string(contentBytes)),
		}

		reqString, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		_, respBody, err := client.DoRequest(http.MethodPut, "/v2/playlists/"+args.ID+"/content", url.Values{}, ptr.Of(string(reqString)), headers)

		return respBody, err
	}))
}

func appendPlaylistContent(playlist *PlaylistResponse, newContent ContentItem) []ContentItem {
	content := []ContentItem{
		newContent,
	}

	if playlist.Content != nil && len(playlist.Content.List) > 0 {
		for _, item := range playlist.Content.List {
			pos, _ := strconv.Atoi(item.Position)
			if pos < newContent.Position {
				content = append([]ContentItem{{
					ID:       item.ID,
					Type:     item.Type,
					Position: pos,
				}}, content...)
			} else {
				content = append(content, ContentItem{
					ID:       item.ID,
					Type:     item.Type,
					Position: pos + 1,
				})
			}
		}
	}

	return content
}

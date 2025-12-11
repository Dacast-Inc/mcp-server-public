package images

import (
	"log/slog"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterUploadSplash(srv *server.MCPServer, client *apiclient.ApiClient) {
	tool := mcp.NewTool("upload_splashscreen_v1",
		mcp.WithDescription("Upload splashscreen image for a given resource (playlist, channel). You must submit either the images as http(s) public URL to the image. File size limit is 10MB, maximum size 5000x5000px. If you cannot send a public link to what the user is attached, do not attempt to find a replacement."),
		mcp.WithInputSchema[UploadImageRequest](),
	)

	srv.AddTool(tool, toolscommon.DacastHandler[UploadImageRequest, UploadImageResponse](func(logger *slog.Logger, args UploadImageRequest, headers map[string]string) ([]byte, error) {
		return uploadImage("splash", client, logger, &args, headers)
	}))
}

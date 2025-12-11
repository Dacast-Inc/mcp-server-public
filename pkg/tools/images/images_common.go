package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools/toolscommon"
	"github.com/google/go-querystring/query"
	"github.com/gotidy/ptr"
	"github.com/mark3labs/mcp-go/server"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type UploadImageRequest struct {
	ID           string `json:"id" jsonschema_description:"ID of the resource to which the image will be uploaded." jsonschema:"required"`
	ResourceType string `json:"resource_type" jsonschema_description:"Type of the resource ('playlist', 'channel')." jsonschema:"required,enum=playlist,enum=channel"`

	ImageURL *string `json:"image_url,omitempty" jsonschema_description:"URL of the image to upload."`
}

type UploadImageResponse struct {
	Status bool `json:"status" jsonschema_description:"Indicates if the image upload was successful."`
}

type PresignedPost struct {
	ACL                 string `json:"acl" url:"acl"`
	Bucket              string `json:"bucket" url:"bucket"`
	Key                 string `json:"key" url:"key"`
	Policy              string `json:"policy" url:"policy"`
	SuccessActionStatus string `json:"success_action_status" url:"success_action_status"`
	Algorithm           string `json:"x-amz-algorithm" url:"x-amz-algorithm"`
	Credential          string `json:"x-amz-credential" url:"x-amz-credential"`
	Date                string `json:"x-amz-date" url:"x-amz-date"`
	Signature           string `json:"x-amz-signature" url:"x-amz-signature"`
	SecurityToken       string `json:"x-amz-security-token,omitempty" url:"x-amz-security-token,omitempty"`
}

func (l UploadImageRequest) Transform() toolscommon.Transformable  { return l }
func (l UploadImageResponse) Transform() toolscommon.Transformable { return l }

func Register(srv *server.MCPServer, client *apiclient.ApiClient) {
	RegisterUploadThumbnail(srv, client)
	RegisterUploadSplash(srv, client)
}

func uploadImage(uploadType string, client *apiclient.ApiClient, logger *slog.Logger, args *UploadImageRequest, headers map[string]string) ([]byte, error) {
	data, err := extractImageFromRequest(args)
	if err != nil {
		return nil, err
	}

	bodyStr := "{\"upload_type\": \"ajax\"}"

	_, body, err := client.DoRequest(http.MethodPost, "/v2/"+args.ResourceType+"/"+args.ID+"/"+uploadType, url.Values{}, ptr.Of(bodyStr), headers)
	if err != nil {
		return nil, err
	}

	var pp PresignedPost
	if err = json.Unmarshal(body, &pp); err != nil {
		return nil, err
	}

	targetUrl := "https://" + pp.Bucket
	values, err := query.Values(pp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query parameters")
	}

	values.Del("bucket")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for key, value := range values {
		err := w.WriteField(key, value[0])
		if err != nil {
			return nil, fmt.Errorf("failed to write field")
		}
	}

	fh, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{fmt.Sprintf(`form-data; name="file"; filename="image.jpg"`)},
		"Content-Type":        []string{"application/octet-stream"},
	})

	if err != nil {
		return nil, err
	}
	if _, err := fh.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, targetUrl, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Error("failed to upload image", "status_code", resp.StatusCode, "body", string(body))

		return nil, fmt.Errorf("failed to upload image to storage: %s", resp.Status)
	}

	return []byte(`{"status": true}`), nil
}

func extractImageFromRequest(args *UploadImageRequest) ([]byte, error) {
	var err error

	var data []byte

	if args.ImageURL != nil && *args.ImageURL != "" {
		resp, err := http.Get(*args.ImageURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10 MB should be enough
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no image data provided")
	}

	if err = validateImage(data, 5000*5000, 5000, 5000); err != nil {
		return nil, err
	}

	return data, nil
}

func validateImage(data []byte, maxPx, maxW, maxH int) error {
	if len(data) == 0 {
		return fmt.Errorf("image data is empty")
	}

	//quick check
	mime := http.DetectContentType(data)
	if !strings.HasPrefix(mime, "image/") {
		return fmt.Errorf("invalid image MIME type: %s", mime)
	}

	//detailed check
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return fmt.Errorf("invalid image data: %w", err)
	}

	if (maxW > 0 && cfg.Width > maxW) ||
		(maxH > 0 && cfg.Height > maxH) ||
		(maxPx > 0 && cfg.Width*cfg.Height > maxPx) {
		return fmt.Errorf("image dimensions exceed allowed size: %dx%d", cfg.Width, cfg.Height)
	}

	return nil
}

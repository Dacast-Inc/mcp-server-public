package toolscommon

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/google/go-querystring/query"
	"github.com/gotidy/ptr"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DacastHandlerFunc[TArgs Transformable] func(logger *slog.Logger, args TArgs, headers map[string]string) ([]byte, error)

func DefaultDacastProxyHandler[TArgs Transformable, TResult Transformable](client *apiclient.ApiClient, method, endpoint string, urlParams []string) server.ToolHandlerFunc {
	return DacastHandler[TArgs, TResult](func(logger *slog.Logger, args TArgs, headers map[string]string) ([]byte, error) {
		*logger = *logger.With("method", method, "rawEndpoint", endpoint)

		q, err := query.Values(args.Transform().(TArgs))
		if err != nil {
			logger.Error("failed to parse query parameters")

			return nil, fmt.Errorf("failed to encode arguments: %v", err)
		}

		finalEndpoint := endpoint
		for _, param := range urlParams {
			if q.Get(param) != "" {
				finalEndpoint = strings.Replace(finalEndpoint, "{"+param+"}", q.Get(param), -1)
			}

			q.Del(param)
		}

		*logger = *logger.With("finalEndpoint", finalEndpoint)

		var bodyString *string
		if method != http.MethodGet && method != http.MethodDelete {
			bodyBytes, err := json.Marshal(args)
			if err != nil {
				logger.Error("failed to marshal request body", "error", err)

				return nil, fmt.Errorf("failed to marshal request body: %v", err)
			}

			bodyString = ptr.Of(string(bodyBytes))
			q = url.Values{}
			headers["Content-Type"] = "application/json"
		}

		bodyStringValue := ""
		if bodyString != nil {
			bodyStringValue = *bodyString
		}
		logger.Info("making API request", "endpoint", endpoint, "query_params", q.Encode(), "body", bodyStringValue)

		_, body, err := client.DoRequest(method, finalEndpoint, q, bodyString, headers)
		if err != nil {
			logger.Error("failed to call api", "error", err, "body", string(body))

			return nil, fmt.Errorf("API request failed: %v", err)
		}

		return body, nil
	})
}

func DacastHandler[TArgs Transformable, TResult Transformable](handler DacastHandlerFunc[TArgs]) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, retErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%v - %s", r, debug.Stack())

				retErr = fmt.Errorf("internal server error")
			}
		}()

		var err error

		logger := slog.With("args", request.GetArguments())

		token, ok := ApiKeyFromContext(ctx)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid token")), nil
		}

		logger = logger.With("token", token)

		headers := map[string]string{
			"X-Api-Key": token,
			"X-Format":  "default",
			"Origin":    "https://mcp.dacast.com",
		}

		ip, _ := ctx.Value("source_ip").(string)

		if ip != "" {
			headers["X-Mcp-Forwarded-For"] = ip

			logger = logger.With("source_ip", ip)
		}

		var args TArgs
		if err = request.BindArguments(&args); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to bind arguments: %v", err)), nil
		}

		body, err := handler(logger, args, headers)
		if err != nil {
			logger.Error("failed to call handler", "error", err)

			return mcp.NewToolResultError(fmt.Sprintf("%v", err)), nil
		}

		toolResult := mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(body),
				},
			},
		}

		var structured TResult
		err = json.Unmarshal(body, &structured)
		if err != nil {
			logger.Error("failed to unmarshal response body", "error", err, "body", string(body))
		} else {
			toolResult.StructuredContent = structured.Transform().(TResult)

			afterTransformBody, err := json.Marshal(toolResult.StructuredContent)
			if err != nil {
				logger.Error("failed to marshal transformed response body", "error", err)
			} else {
				toolResult.Content = []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: string(afterTransformBody),
					},
				}
			}
		}

		return &toolResult, nil
	}
}

func ApiKeyFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value("token").(string)
	return v, ok && v != ""
}

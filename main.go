package main

import (
	"context"
	"os"
	"time"

	"github.com/Dacast-Inc/mcp-server-public/pkg/apiclient"
	"github.com/Dacast-Inc/mcp-server-public/pkg/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	apiToken := os.Getenv("DACAST_API_TOKEN")
	if apiToken == "" {
		panic("DACAST_API_TOKEN environment variable is not set")
	}

	apiClient := apiclient.NewApiClient("https://developer.dacast.com", 60*time.Second)

	srv := server.NewMCPServer(
		"Dacast.com MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	tools.Register(srv, apiClient)

	server.ServeStdio(srv, withAPIToken(apiToken))
}

func withAPIToken(token string) server.StdioOption {
	return server.WithStdioContextFunc(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "token", token)
	})
}

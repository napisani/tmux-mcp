package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/napisani/tmux-mcp/internal/resource"
	"github.com/napisani/tmux-mcp/internal/tool"
)

func setupLogging(verbose bool) {
	if verbose {
		logFile, err := os.OpenFile("/tmp/tmux-mcp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		log.SetOutput(logFile)
		log.Println("Logging started")
	} else {
		log.SetOutput(io.Discard)
	}
}

func main() {

	setupLogging(true)

	// Create a new MCP server
	s := server.NewMCPServer(
		"tmux-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	// // Define a simple tool
	// tool := mcp.NewTool("hello_world",
	// 	mcp.WithDescription("Say hello to someone"),
	// 	mcp.WithString("name",
	// 		mcp.Required(),
	// 		mcp.Description("Name of the person to greet"),
	// 	),
	// )

	s.AddResource(resource.GetSessionsListResource(), resource.HandleSessionsList)
	s.AddResource(resource.GetWindowsListResource(), resource.HandleWindowsList)
	s.AddResource(resource.GetPanesListResource(), resource.HandlePanesList)
	s.AddResource(resource.GetPaneOutputResource(), resource.HandleGetPaneOutput)

	// Add tool handler
	// s.AddTool(tool, helloHandler)
	s.AddTool(tool.GetSessionListTool(), tool.HandleSessionsList)
	s.AddTool(tool.GetWindowsListTool(), tool.HandleWindowsList)
	s.AddTool(tool.GetPaneOutputTool(), tool.HandleGetPaneOutput)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s! 👋", name)), nil
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/napisani/tmux-mcp/internal/tmuxmcp"
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

	s.AddResource(tmuxmcp.GetSessionsListResource(), tmuxmcp.HandleSessionsListResource)
	s.AddResource(tmuxmcp.GetWindowsListResource(), tmuxmcp.HandleWindowsListResource)
	s.AddResource(tmuxmcp.GetPanesListResource(), tmuxmcp.HandlePanesListResource)
	s.AddResourceTemplate(tmuxmcp.GetPaneOutputResourceTemplate(), tmuxmcp.HandlePaneOutputResourceTemplate)

	// Add tool handler
	// s.AddTool(tool, helloHandler)
	s.AddTool(tmuxmcp.GetSessionListTool(), tmuxmcp.HandleSessionsList)
	s.AddTool(tmuxmcp.GetWindowsListTool(), tmuxmcp.HandleWindowsList)
	s.AddTool(tmuxmcp.GetPaneOutputTool(), tmuxmcp.HandleGetPaneOutput)
	s.AddTool(tmuxmcp.GetPaneByTagTool(), tmuxmcp.HandleFindPaneByTag)

	// Execute command tool
	s.AddTool(tmuxmcp.ExecuteCommandTool(), tmuxmcp.HandleExecuteCommand)

	// Execute async tool
	s.AddTool(tmuxmcp.ExecuteCommandAsyncTool(), tmuxmcp.HandleExecuteAsync)

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

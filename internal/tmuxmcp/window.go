package tmuxmcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

// Tool definitions

func GetWindowsListTool() mcp.Tool {
	return mcp.NewTool("list_windows",
		mcp.WithDescription("List all of the tmux windows in the current tmux server"),
	)
}

func HandleWindowsList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	windows, err := getWindowsList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get windows: %v", err)), nil
	}

	jsonContent, err := json.Marshal(windows)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal windows: %v", err)), nil
	}

	content := fmt.Sprintf("%s \n Windows in the tmux server:\n```json\n%s```", desc.WindowListDescription, string(jsonContent))

	return mcp.NewToolResultText(content), nil
}

// Resource definitions

func GetWindowsListResource() mcp.Resource {
	description := desc.WindowListDescription

	return mcp.NewResource(
		"windows://",
		"Tmux Windows",
		mcp.WithResourceDescription(description),
		mcp.WithMIMEType("application/json"),
	)
}

func HandleWindowsListResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	windows, err := getWindowsList()
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(windows)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(configJSON),
		},
	}, nil
}

// Shared implementation

func getWindowsList() ([]*gotmux.Window, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tmux: %w", err)
	}
	
	windows, err := tmux.ListAllWindows()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}
	
	return windows, nil
}

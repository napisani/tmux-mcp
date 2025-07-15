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

func GetSessionListTool() mcp.Tool {
	return mcp.NewTool("list_sessions",
		mcp.WithDescription("List all of the tmux sessions in the current tmux server"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
}

func HandleSessionsList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessions, err := getSessionsList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get sessions: %v", err)), nil
	}

	jsonContent, err := json.Marshal(sessions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal sessions: %v", err)), nil
	}

	content := fmt.Sprintf("%s \n Sessions in the tmux server:\n```json\n%s```", desc.SessionListDescription, string(jsonContent))

	return mcp.NewToolResultText(content), nil
}

// Resource definitions

func GetSessionsListResource() mcp.Resource {
	description := desc.SessionListDescription 
	return mcp.NewResource(
		"sessions://",
		"Tmux Sessions",
		mcp.WithResourceDescription(description),
		mcp.WithMIMEType("application/json"),
	)
}

func HandleSessionsListResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	sessions, err := getSessionsList()
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(sessions)
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

func getSessionsList() ([]*gotmux.Session, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tmux: %w", err)
	}
	
	sessions, err := tmux.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	
	return sessions, nil
}

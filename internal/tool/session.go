package tool

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

// Define a simple tool
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
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to tmux: %v", err)), nil
	}
	sessions, err := tmux.ListSessions()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list sessions: %v", err)), nil
	}
	jsonContent, err := json.Marshal(sessions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal sessions: %v", err)), nil
	}

	content := fmt.Sprintf("%s \n Sessions in the tmux server:\n```json\n%s```", desc.SessionListDescription, string(jsonContent))

	return mcp.NewToolResultText(content), nil

}

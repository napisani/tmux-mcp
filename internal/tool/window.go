package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

func GetWindowsListTool() mcp.Tool {
	return mcp.NewTool("list_windows",
		mcp.WithDescription("List all of the tmux windows in the current tmux server"),
	)
}

func HandleWindowsList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to tmux: %v", err)), nil
	}
	windows, err := tmux.ListAllWindows()

	jsonContent, err := json.Marshal(windows)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal windows: %v", err)), nil
	}

	content := fmt.Sprintf("%s \n Windows in the tmux server:\n```json\n%s```", desc.WindowListDescription, string(jsonContent))

	return mcp.NewToolResultText(content), nil
}

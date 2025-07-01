// internal/tool/output.go
package tool

import (
	"context"
	"fmt"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
)

// Define a tool for getting pane output
func GetPaneOutputTool() mcp.Tool {
	return mcp.NewTool("get_pane_output",
		mcp.WithDescription("Get the output content from a specific tmux pane"),
		mcp.WithString("pane_id",
			mcp.Required(),
			mcp.Description("ID of the pane to get output from"),
		),
		mcp.WithNumber("lines",
			mcp.Description("Number of lines to capture from the pane history (default: 100)"),
		),
	)
}

func HandleGetPaneOutput(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to tmux: %v", err)), nil
	}

	paneId, err := request.RequireString("pane_id")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get pane_id: %v", err)), nil
	}

	// Default to 100 lines if not specified
	lines := request.GetInt("lines", 100)

	// Verify the pane exists
	panes, err := tmux.ListAllPanes()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list panes: %v", err)), nil
	}

	var paneExists bool
	for _, pane := range panes {
		if pane.Id == paneId {
			paneExists = true
			break
		}
	}

	if !paneExists {
		return mcp.NewToolResultError(fmt.Sprintf("Pane with ID %s not found", paneId)), nil
	}

	// Capture pane content using tmux capture-pane command
	output, err := tmux.Command(fmt.Sprintf("capture-pane -p -t %s -S -%d", paneId, lines))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to capture pane output: %v", err)), nil
	}

	content := fmt.Sprintf("Output from pane %s (last %d lines):\n```\n%s\n```", paneId, lines, output)

	return mcp.NewToolResultText(content), nil
}

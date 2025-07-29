package tmuxmcp

import (
	"context"
	"fmt"
	"log"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
)

// ExecuteAsyncTool defines the MCP tool that kicks off a command in a given
// pane without waiting to capture its output. Think of it as a fire-and-forget
// variant of ExecuteCommand.
func ExecuteCommandAsyncTool() mcp.Tool {
	return mcp.NewTool("execute_command_async",
		mcp.WithDescription("Execute a shell command in a tmux pane without waiting for output (async)"),
		mcp.WithString("pane_id",
			mcp.Required(),
			mcp.DefaultString("0"),
			mcp.Description("ID of the pane to run the command in"),
		),
		mcp.WithNumber("window_index",
			mcp.Required(),
			mcp.DefaultNumber(0),
			mcp.Description("Index of the window containing the pane"),
		),
		mcp.WithString("session_name",
			mcp.Required(),
			mcp.Description("Name of the tmux session containing the pane"),
		),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Shell command to execute asynchronously"),
		),
	)
}

func HandleExecuteAsync(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paneID, err := req.RequireString("pane_id")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get pane_id: %v", err)), nil
	}
	windowIndex, err := req.RequireInt("window_index")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get window_index: %v", err)), nil
	}
	sessionName, err := req.RequireString("session_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get session_name: %v", err)), nil
	}
	cmd, err := req.RequireString("command")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get command: %v", err)), nil
	}

	if err := sendCommandToPane(sessionName, windowIndex, paneID, cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("command dispatched successfully"), nil
}

func sendCommandToPane(session string, winIdx int, paneID, command string) error {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return fmt.Errorf("failed to connect to tmux: %w", err)
	}

	target := fmt.Sprintf("%s:%d.%s", session, winIdx, paneID)
	log.Printf("[execute_async] sending command to pane %s: %s", target, command)

	if _, err := tmux.Command("send-keys", "-t", target, command, "C-m"); err != nil {
		return fmt.Errorf("failed to send command to pane: %w", err)
	}
	return nil
}

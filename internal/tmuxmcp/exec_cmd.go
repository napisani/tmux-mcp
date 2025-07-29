package tmuxmcp

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
)

// ExecuteCommandTool returns the MCP tool definition for executing a shell
// command inside a specific tmux pane while tee'ing output to a temporary
// file which is returned to the caller.
func ExecuteCommandTool() mcp.Tool {
	return mcp.NewTool("execute_command",
		mcp.WithDescription("Execute a shell command in a tmux pane and return combined stdout/stderr"),
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
			mcp.Description("Shell command to execute"),
		),
	)
}

// HandleExecuteCommand executes the requested command in the target pane and
// returns its output (stdout & stderr merged) back to the caller.
func HandleExecuteCommand(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	userCmd, err := req.RequireString("command")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get command: %v", err)), nil
	}

	output, execErr := executeCommandInPane(sessionName, windowIndex, paneID, userCmd)
	if execErr != nil {
		return mcp.NewToolResultError(execErr.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

// executeCommandInPane sends a shell command wrapped with tee to the target
// pane, waits briefly, reads the captured output and returns it.
func executeCommandInPane(session string, windowIdx int, paneID, cmd string) (string, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return "", fmt.Errorf("failed to connect to tmux: %w", err)
	}

	// Verify pane exists to return meaningful error early.
	panes, err := tmux.ListAllPanes()
	if err != nil {
		return "", fmt.Errorf("failed to list panes: %w", err)
	}
	var found bool
	for _, p := range panes {
		if matchPaneId(paneID, &McpPane{Pane: p}) {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("pane with ID %s not found", paneID)
	}

	// Create temp file for tee output.
	tmpFile, err := os.CreateTemp("", "tmux-mcp-exec-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Unique wait token so concurrent calls don't clash.
	token := fmt.Sprintf("mcp_%d", time.Now().UnixNano())

	// Build the shell wrapper:
	//   (<cmd> 2>&1 | tee tmp); tmux wait-for -S <token>
	wrapped := fmt.Sprintf("bash -c '(%s 2>&1 | tee %s); tmux wait-for -S %s'", escapeSingleQuotes(cmd), tmpPath, token)

	target := fmt.Sprintf("%s:%d.%s", session, windowIdx, paneID)
	log.Printf("Sending command to pane %s: %s", target, cmd)

	// Inject command.
	if _, err := tmux.Command("send-keys", "-t", target, wrapped, "C-m"); err != nil {
		return "", fmt.Errorf("failed to send command to pane: %w", err)
	}

	// Block until the pane-side command signals completion.
	if _, err := tmux.Command("wait-for", token); err != nil {
		return "", fmt.Errorf("wait-for failure: %w", err)
	}

	// Read captured output after command has finished.
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read captured output: %w", err)
	}
	os.Remove(tmpPath)

	return string(data), nil
}

// escapeSingleQuotes makes a string safe to embed within single-quoted bash
// string by closing, escaping, and reopening the quote.
func escapeSingleQuotes(s string) string {
	// Replace ' with '\''
	return strings.ReplaceAll(s, "'", "'\\''")
}

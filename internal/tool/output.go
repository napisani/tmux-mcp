// internal/tool/output.go
package tool

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/napisani/tmux-mcp/internal/exec"
)

func matchPaneId(paneId string, pane *gotmux.Pane) bool {
	return paneId == pane.Id || string('%')+paneId == pane.Id
}

// Define a tool for getting pane output
func GetPaneOutputTool() mcp.Tool {
	return mcp.NewTool("get_pane_output",
		mcp.WithDescription("Get the output content from a specific tmux pane"),
		mcp.WithString("pane_id",
			mcp.Required(),
			mcp.Description("ID of the pane to get output from"),
		),
		mcp.WithNumber("lines",
			mcp.Required(),
			mcp.Description("Number of lines to capture from the pane history (default: 100)"),
		),
	)
}

func HandleGetPaneOutput(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		log.Printf("Error: Failed to connect to tmux: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to tmux: %v", err)), nil
	}

	paneId, err := request.RequireString("pane_id")
	log.Printf("Looking for Pane ID: %s", paneId)
	if err != nil {
		log.Printf("Error: Failed to get pane_id: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get pane_id: %v", err)), nil
	}

	// Default to 100 lines if not specified
	lines := request.GetInt("lines", 100)

	// Verify the pane exists
	panes, err := tmux.ListAllPanes()
	for _, pane := range panes {
		// print the id
		log.Printf("Pane ID: %s", pane.Id)
	}
	if err != nil {
		log.Printf("Error: Failed to list panes: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list panes: %v", err)), nil
	}

	var paneExists bool
	for _, pane := range panes {
		if matchPaneId(paneId, pane) {
			paneExists = true
			break
		}
	}

	if !paneExists {
		log.Printf("Error: Pane with ID %s not found", paneId)
		return mcp.NewToolResultError(fmt.Sprintf("Pane with ID %s not found", paneId)), nil
	}

	// Create a temporary file to store the pane content
	tmpFile, err := os.CreateTemp("", "tmux-mcp-capture-*.txt")
	if err != nil {
		log.Printf("Error: Failed to create temporary file: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create temporary file: %v", err)), nil
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file when done

	// Capture pane content using tmux capture-pane command and write to the temporary file
	//capture-pane -S -32768 ; save-buffer %1 ; delete-buffer
	cmd := fmt.Sprintf("capture-pane -t %s -S -%d \\; save-buffer %s ", paneId, lines, tmpFile.Name())
	cmder := exec.NewCommander()
	_, err = cmder.RunCommand("", cmd)
	if err != nil {
		log.Printf("Error: Failed to capture pane output: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to capture pane output: %v", err)), nil
	}

	// Read the file content
	output, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Printf("Error: Failed to read captured content from file: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read captured content: %v", err)), nil
	}
	content := fmt.Sprintf("Output from pane %s (last %d lines):\n```\n%s\n```", paneId, lines, output)

	return mcp.NewToolResultText(content), nil
}

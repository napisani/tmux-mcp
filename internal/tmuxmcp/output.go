package tmuxmcp

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
)

const defaultLines = 200

func matchPaneId(paneId string, pane *McpPane) bool {
	return paneId == pane.Id || string('%')+paneId == pane.Id
}

// Tool definitions
func GetPaneOutputTool() mcp.Tool {
	return mcp.NewTool("get_pane_output",
		mcp.WithDescription("Get the output content from a specific tmux pane"),
		mcp.WithString("pane_id",
			mcp.Required(),
			mcp.DefaultString("0"),
			mcp.Description("ID of the pane to get output from"),
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
		mcp.WithNumber("lines",
			mcp.DefaultNumber(defaultLines),
			mcp.Description(fmt.Sprintf("Number of lines to capture from the pane history (default: %d)", defaultLines)),
		),
	)
}

func HandleGetPaneOutput(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paneId, err := request.RequireString("pane_id")
	if err != nil {
		log.Printf("Error: Failed to get pane_id: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get pane_id: %v", err)), nil
	}

	// Default to 100 lines if not specified
	lines := request.GetInt("lines", defaultLines)

	sessionName, err := request.RequireString("session_name")
	if err != nil {
		log.Printf("Error: Failed to get session_name: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get session_name: %v", err)), nil
	}

	windowIndex, err := request.RequireInt("window_index")
	if err != nil {
		log.Printf("Error: Failed to get window_index: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get window_index: %v", err)), nil
	}

	content, err := getPaneOutput(paneId, sessionName, windowIndex, lines)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get pane output: %v", err)), nil
	}

	return mcp.NewToolResultText(content), nil
}

// // Resource definitions
// func GetPaneOutputResource() mcp.Resource {
// 	description := desc.PaneOutputDescription
// 	return mcp.NewResource(
// 		"pane_output://{session_name}/{window_index}/{pane_id}",
// 		"Tmux Pane Output",
// 		mcp.WithResourceDescription(description),
// 		mcp.WithMIMEType("text/markdown"),
// 	)
// }

// func HandlePaneOutputResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
// 	log.Printf("Handling pane output resource request: %s", request.Params.URI)
// 	fields := request.Params.URI[len("pane_output://"):]
// 	var sessionName, paneId string
// 	var windowIndex int
// 	_, err := fmt.Sscanf(fields, "%s/%d/%s", &sessionName, &windowIndex, &paneId)
// 	if err != nil {
// 		log.Printf("Error: Failed to parse URI: %v", err)
// 		return nil, err
// 	}

// 	lines := defaultLines
// 	content, err := getPaneOutput(paneId, sessionName, windowIndex, lines)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []mcp.ResourceContents{
// 		mcp.TextResourceContents{
// 			URI:      request.Params.URI,
// 			MIMEType: "text/markdown",
// 			Text:     content,
// 		},
// 	}, nil
// }

// Shared implementation
func getPaneOutput(paneId string, sessionName string, windowIndex int, lines int) (string, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		log.Printf("Error: Failed to connect to tmux: %v", err)
		return "", fmt.Errorf("failed to connect to tmux: %w", err)
	}

	// Verify the pane exists
	panes, err := tmux.ListAllPanes()
	if err != nil {
		log.Printf("Error: Failed to list panes: %v", err)
		return "", fmt.Errorf("failed to list panes: %w", err)
	}

	var paneExists bool
	for _, pane := range panes {
		if matchPaneId(paneId, &McpPane{Pane: pane}) {
			paneExists = true
			break
		}
	}

	if !paneExists {
		log.Printf("Error: Pane with ID %s not found", paneId)
		return "", fmt.Errorf("pane with ID %s not found", paneId)
	}

	// Create a temporary file to store the pane content
	tmpFile, err := os.CreateTemp("", "tmux-mcp-capture-*.txt")
	if err != nil {
		log.Printf("Error: Failed to create temporary file: %v", err)
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file when done

	cmd := []string{
		"capture-pane",
		"-t", fmt.Sprintf("%s:%d.%s", sessionName, windowIndex, paneId),
		"-S", fmt.Sprintf("-%d", lines),
		";",
		"save-buffer", tmpFile.Name(),
	}

	log.Printf("Command to capture pane output: %s", cmd)
	log.Printf("Session %s", sessionName)
	_, err = tmux.Command(cmd...)
	if err != nil {
		log.Printf("Error: Failed to capture pane output: %v", err)
		return "", fmt.Errorf("failed to capture pane output: %w", err)
	}

	// Read the file content
	output, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Printf("Error: Failed to read captured content from file: %v", err)
		return "", fmt.Errorf("failed to read captured content: %w", err)
	}

	return fmt.Sprintf("Output from pane %s (last %d lines):\n```\n%s\n```", paneId, lines, output), nil
}

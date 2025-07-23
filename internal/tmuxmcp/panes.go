package tmuxmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

// Resource definitions

func GetPanesListResource() mcp.Resource {
	description := desc.PaneOutputDescription
	return mcp.NewResource(
		"panes://",
		"Tmux Panes",
		mcp.WithResourceDescription(description),
		mcp.WithMIMEType("application/json"),
	)
}

func HandlePanesListResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	panes, err := getPanesList()
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(panes)
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

func getMcpTag(tmux *gotmux.Tmux, p *gotmux.Pane) string {
	if p == nil {
		return ""
	}
	out, err := tmux.Command("show-option", "-p", "-v", "-t", p.Id, "@mcptag")
	if err != nil {
		// Silently ignore error and return empty tag
		return ""
	}
	return strings.TrimSpace(out)
}

func getPanesList() ([]*McpPane, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tmux: %w", err)
	}

	panes, err := tmux.ListAllPanes()
	if err != nil {
		return nil, fmt.Errorf("failed to list panes: %w", err)
	}

	wrapped := make([]*McpPane, len(panes))
	for i, pane := range panes {
		tag := getMcpTag(tmux, pane)
		log.Printf("Pane %d: %s, McpTag: %s", i, pane.Id, tag)
		wrapped[i] = &McpPane{Pane: pane, McpTag: tag}
		log.Println(pane.SessionName)
	}

	return wrapped, nil
}

// Helper performing the fuzzy-matching search on an in-memory slice so it can be unit-tested.
// Matching rule: case-insensitive substring containment.
// Returns error if 0 or >1 panes match.
func findPaneByTag(tag string, panes []*McpPane) (*McpPane, error) {
	var matches []*McpPane
	lo := strings.ToLower(tag)
	for _, p := range panes {
		if strings.Contains(strings.ToLower(p.McpTag), lo) {
			matches = append(matches, p)
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no pane found with tag %q", tag)
	case 1:
		return matches[0], nil
	default:
		return nil, fmt.Errorf("multiple panes (%d) match tag %q", len(matches), tag)
	}
}

// Tool definitions
func GetPaneByTagTool() mcp.Tool {
	return mcp.NewTool("find_pane_by_tag",
		mcp.WithDescription("Find a tmux pane by fuzzy-matching its McpTag"),
		mcp.WithString("tag",
			mcp.Required(),
			mcp.Description("Tag to fuzzy-match against pane McpTag"),
		),
	)
}

func HandleFindPaneByTag(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag, err := request.RequireString("tag")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get tag: %v", err)), nil
	}

	panes, err := getPanesList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list panes: %v", err)), nil
	}

	pane, err := findPaneByTag(tag, panes)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bytes, _ := json.Marshal(pane)
	return mcp.NewToolResultText(string(bytes)), nil
}

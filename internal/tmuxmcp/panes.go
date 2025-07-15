package tmuxmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func getPanesList() ([]*gotmux.Pane, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tmux: %w", err)
	}

	panes, err := tmux.ListAllPanes()
	if err != nil {
		return nil, fmt.Errorf("failed to list panes: %w", err)
	}

	for _, pane := range panes {
		log.Println(pane.SessionName)

	}

	return panes, nil
}

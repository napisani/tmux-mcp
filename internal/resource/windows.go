package resource

import (
	"context"
	"encoding/json"
	"log"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

func GetWindowsListResource() mcp.Resource {
	description := desc.WindowListDescription

	return mcp.NewResource(
		"windows://",
		"Tmux Windows",
		mcp.WithResourceDescription(description),
		mcp.WithMIMEType("application/json"),
	)
}

func HandleWindowsList(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	windows, err := tmux.ListAllWindows()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	configJSON, err := json.Marshal(windows)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(configJSON),
			},
		},
		nil
}

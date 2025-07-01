package resource

import (
	"context"
	"fmt"
	"log"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/napisani/tmux-mcp/internal/desc"
)

func GetPaneOutputResource() mcp.Resource {
	description := desc.PaneOutputDescription
	return mcp.NewResource(
		"pane_output://{session_id}/{window_id}/{pane_id}",
		"Tmux Pane Output",
		mcp.WithResourceDescription(description),
		mcp.WithMIMEType("text/markdown"),
	)
}

func HandleGetPaneOutput(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fields := request.Params.URI[len("pane_output://"):]
	var sessionId, windowId, paneId string
	_, err = fmt.Sscanf(fields, "%s/%s/%s", &sessionId, &windowId, &paneId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	lines := 100
	output, err := tmux.Command(fmt.Sprintf("capture-pane -p -t %s -S -%d", paneId, lines))
	content := fmt.Sprintf("Output from pane %s (last %d lines):\n```\n%s\n```", paneId, lines, output)

	return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		},
		nil
}

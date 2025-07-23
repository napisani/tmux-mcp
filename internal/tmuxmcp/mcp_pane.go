package tmuxmcp

import "github.com/GianlucaP106/gotmux/gotmux"

// McpPane wraps a gotmux.Pane with additional MCP-specific metadata.
// It embeds the original *gotmux.Pane so all its fields and methods are retained.
// McpTag can be used by higher-level tools to attach arbitrary labels to a pane.
//
// NOTE: All references to gotmux.Pane within this code-base should prefer McpPane
// to allow future MCP-specific extensions without touching call-sites.
type McpPane struct {
	*gotmux.Pane `json:""`
	McpTag       string `json:"McpTag,omitempty"`
}

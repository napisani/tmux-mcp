package tmuxmcp

import (
	"github.com/napisani/tmux-mcp/internal/testutil"
	"testing"
)

func TestGetPanesList(t *testing.T) {
	session := testutil.CreateTestSession(t)

	panes, err := getPanesList()
	if err != nil {
		t.Fatalf("getPanesList returned error: %v", err)
	}
	found := false
	for _, p := range panes {
		if p.SessionName == session {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("pane for session %s not found", session)
	}
}


func TestFindPaneByTag(t *testing.T) {
	panes := []*McpPane{
		{McpTag: "build"},
		{McpTag: "deploy"},
		{McpTag: "builder"},
	}

	// Single unique fuzzy match
	p, err := findPaneByTag("depl", panes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.McpTag != "deploy" {
		t.Fatalf("expected 'deploy', got %s", p.McpTag)
	}
	// Ambiguous match
	if _, err := findPaneByTag("build", panes); err == nil {
		t.Fatalf("expected error for ambiguous match, got nil")
	}

	// No match
	if _, err := findPaneByTag("xyz", panes); err == nil {
		t.Fatalf("expected error for no match, got nil")
	}
}

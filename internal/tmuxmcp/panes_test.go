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

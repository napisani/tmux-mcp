package tmuxmcp

import (
	"github.com/napisani/tmux-mcp/internal/testutil"
	"testing"
)

func TestGetPaneOutput(t *testing.T) {
	session := testutil.CreateTestSession(t)

	panes, err := getPanesList()
	if err != nil {
		t.Fatalf("could not list panes: %v", err)
	}
	var paneId string
	for _, p := range panes {
		if p.SessionName == session {
			paneId = p.Id
			break
		}
	}
	if paneId == "" {
		t.Fatalf("no pane found for session %s", session)
	}

	output, err := getPaneOutput(paneId, session, 0, 10)
	if err != nil {
		t.Fatalf("getPaneOutput error: %v", err)
	}
	if output == "" {
		t.Fatalf("expected non-empty pane output")
	}
}

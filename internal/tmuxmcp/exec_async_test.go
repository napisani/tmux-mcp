package tmuxmcp

import (
	"testing"

	"github.com/napisani/tmux-mcp/internal/testutil"
)

func TestExecuteAsync(t *testing.T) {
	session := testutil.CreateTestSession(t)
	panes, err := getPanesList()
	if err != nil {
		t.Fatalf("list panes: %v", err)
	}
	var paneID string
	for _, p := range panes {
		if p.SessionName == session {
			paneID = p.Id
			break
		}
	}
	if paneID == "" {
		t.Fatalf("no pane found in session %s", session)
	}

	// Should not error; we can't easily assert output.
	if err := sendCommandToPane(session, 0, paneID, "echo async_works"); err != nil {
		t.Fatalf("sendCommandToPane: %v", err)
	}
}

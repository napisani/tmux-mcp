package tmuxmcp

import (
	"strings"
	"testing"

	"github.com/napisani/tmux-mcp/internal/testutil"
)

func TestExecuteCommandInPane(t *testing.T) {
	session := testutil.CreateTestSession(t)

	panes, err := getPanesList()
	if err != nil {
		t.Fatalf("could not list panes: %v", err)
	}
	var paneID string
	for _, p := range panes {
		if p.SessionName == session {
			paneID = p.Id
			break
		}
	}
	if paneID == "" {
		t.Fatalf("no pane found for session %s", session)
	}

	out, err := executeCommandInPane(session, 0, paneID, "echo hello_mcp")
	if err != nil {
		t.Fatalf("executeCommandInPane error: %v", err)
	}
	if len(out) == 0 {
		t.Fatalf("expected non-empty output")
	}
	if !contains(out, "hello_mcp") {
		t.Fatalf("output did not contain expected text: %q", out)
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

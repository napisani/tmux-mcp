package tmuxmcp

import (
	"testing"

	"github.com/napisani/tmux-mcp/internal/testutil"
)

func TestGetSessionsList(t *testing.T) {
	session := testutil.CreateTestSession(t)

	sessions, err := getSessionsList()
	if err != nil {
		t.Fatalf("getSessionsList returned error: %v", err)
	}
	found := false
	for _, s := range sessions {
		if s.Name == session {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("session %s not found in sessions list", session)
	}
}

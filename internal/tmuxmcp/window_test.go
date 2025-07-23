package tmuxmcp

import (
	"github.com/napisani/tmux-mcp/internal/testutil"
	"testing"
)

func TestGetWindowsList(t *testing.T) {
	_ = testutil.CreateTestSession(t)

	windows, err := getWindowsList()
	if err != nil {
		t.Fatalf("getWindowsList returned error: %v", err)
	}
	if len(windows) == 0 {
		t.Fatalf("expected at least one window, got 0")
	}
}

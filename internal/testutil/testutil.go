package testutil

import (
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
)

// CreateTestSession spins up a detached tmux session that tests can
// interact with. It returns the session name and registers a cleanup
// callback to kill the session once the test finishes.
func CreateTestSession(t *testing.T) string {
	t.Helper()
	session := "mcp_test_" + uuid.NewString()

	cmd := exec.Command("tmux", "new-session", "-d", "-s", session, "sh")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to create test tmux session: %v – %s", err, out)
	}

	// Give tmux a moment to fully initialize the session.
	time.Sleep(200 * time.Millisecond)

	t.Cleanup(func() {
		exec.Command("tmux", "kill-session", "-t", session).Run()
	})

	return session
}

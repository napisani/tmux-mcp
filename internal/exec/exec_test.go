package exec

import (
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
)

// createTestSession spins up a detached tmux session that the tests can
// interact with.  It returns the session name and registers a Cleanup
// callback to kill the session once the test finishes.
func createTestSession(t *testing.T) string {
	t.Helper()
	session := "mcp_test_" + uuid.NewString()

	// Start a new detached session running a long-lived shell
	cmd := exec.Command("tmux", "new-session", "-d", "-s", session, "sh")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to create test tmux session: %v – %s", err, out)
	}

	// Give tmux a moment to fully create the session – this avoids flaky tests
	time.Sleep(200 * time.Millisecond)

	// Register cleanup to kill the session
	t.Cleanup(func() {
		exec.Command("tmux", "kill-session", "-t", session).Run()
	})

	return session
}

func TestCommander_Run(t *testing.T) {
	// Ensure we have a live session that list-sessions can return
	session := createTestSession(t)
	if session == "" {
		t.Fatal("empty session name returned")
	}

	c := NewCommander()
	out, err := c.Run("list-sessions", "-F", "#S")
	if err != nil {
		t.Fatalf("Commander.Run returned error: %v", err)
	}

	// The output is a newline-separated list of session names – make sure our
	// session is present.
	if !containsLine(out, session) {
		t.Fatalf("expected session %s to be listed, got: %q", session, out)
	}
}

// containsLine reports whether the given multiline string contains exactly the
// supplied line.
func containsLine(haystack, needle string) bool {
	for _, l := range splitLines(haystack) {
		if l == needle {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\n' || c == '\r' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(c)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

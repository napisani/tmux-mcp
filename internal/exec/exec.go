package exec

import (
	"fmt"
	"os/exec"
	"strings"
)

// Commander provides functionality to interact with tmux sessions
type Commander struct {
	TmuxPath string // Path to tmux binary
}

// NewCommander creates a new Commander with the default tmux path
func NewCommander() *Commander {
	return &Commander{
		TmuxPath: "tmux", // Default assumes tmux is in PATH
	}
}

// RunCommand executes a command in the specified tmux session.
// If sessionID is empty, the command runs in the default/current session.
// Returns stdout output and any error encountered.
func (c *Commander) RunCommand(sessionID string, command string) (string, error) {
	// Validate command input
	if command == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	var cmd *exec.Cmd

	// If sessionID is provided, run in that specific session
	// Otherwise, run in the default session
	if sessionID != "" {
		cmd = exec.Command(c.TmuxPath, "run-shell", "-t", sessionID, command)
	} else {
		cmd = exec.Command(c.TmuxPath, "run-shell", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command in tmux session: %w, output: %s",
			err, strings.TrimSpace(string(output)))
	}

	return strings.TrimSpace(string(output)), nil
}

// SessionExists checks if the specified tmux session exists
// If sessionID is empty, it checks if any tmux session exists
func (c *Commander) SessionExists(sessionID string) (bool, error) {
	var cmd *exec.Cmd

	if sessionID != "" {
		cmd = exec.Command(c.TmuxPath, "has-session", "-t", sessionID)
	} else {
		// Check if there's any session running
		cmd = exec.Command(c.TmuxPath, "has-session")
	}

	err := cmd.Run()
	if err != nil {
		// Exit code 1 typically means session doesn't exist
		return false, nil
	}

	return true, nil
}

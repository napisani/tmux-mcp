package exec

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Commander struct {
	// Optional path to tmux binary, defaults to "tmux" in PATH
	TmuxPath string
}

func NewCommander() *Commander {
	return &Commander{
		TmuxPath: "tmux",
	}
}

// Run executes a tmux command and returns its output
func (c *Commander) Run(args ...string) (string, error) {
	cmd := exec.Command(c.TmuxPath, args...)
	log.Println("Executing tmux command:", cmd.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("tmux command failed: %v, stderr: %s", err, stderr.String())
		return "", fmt.Errorf("tmux command failed: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// RunRaw executes a raw tmux command string
func (c *Commander) RunRaw(cmdStr string) (string, error) {
	args := []string{}
	if cmdStr != "" {
		args = append(args, strings.Split(cmdStr, " ")...)
	}
	return c.Run(args...)
}

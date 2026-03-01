// Package services provides an abstraction over shell command execution.
package services

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/chmenegatti/devdash/internal/logs"
)

// CommandResult holds the stdout, stderr, and error from a command.
type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

// RunCommand executes a shell command in the given directory and returns the
// combined result. This thin wrapper enables future mocking and testability.
func RunCommand(ctx context.Context, dir, name string, args ...string) CommandResult {
	started := time.Now()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logs.Errorf(
			"command failed after %s | dir=%s | cmd=%s %s | err=%v | stdout=%s | stderr=%s",
			time.Since(started).Round(time.Millisecond),
			dir,
			name,
			strings.Join(args, " "),
			err,
			trimForLog(stdout.String()),
			trimForLog(stderr.String()),
		)
	}

	return CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

func trimForLog(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= 500 {
		return s
	}
	return s[:500] + "..."
}

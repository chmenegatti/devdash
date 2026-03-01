// Package services provides an abstraction over shell command execution.
package services

import (
	"bytes"
	"context"
	"os/exec"
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
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

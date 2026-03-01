package modules

import (
	"context"
	"strings"
	"time"

	"github.com/chmenegatti/devdash/internal/services"
	"github.com/chmenegatti/devdash/internal/state"
)

// RunGitStatus executes `git status --short` and categorizes changed files.
// Blocking — call from a tea.Cmd.
func RunGitStatus(projectDir string) state.GitResult {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "git", "status", "--short")

	return parseGitOutput(res)
}

// parseGitOutput categorizes `git status --short` lines by their prefix.
func parseGitOutput(res services.CommandResult) state.GitResult {
	if res.Err != nil {
		return state.GitResult{
			Status: state.StatusError,
			Err:    res.Err.Error(),
		}
	}

	var modified, added, deleted, other []string

	for _, line := range strings.Split(res.Stdout, "\n") {
		if len(line) < 3 {
			continue
		}
		prefix := line[:2]
		file := strings.TrimSpace(line[2:])

		switch {
		case strings.Contains(prefix, "M"):
			modified = append(modified, file)
		case strings.Contains(prefix, "A"):
			added = append(added, file)
		case strings.Contains(prefix, "D"):
			deleted = append(deleted, file)
		case strings.Contains(prefix, "?"):
			other = append(other, file)
		default:
			other = append(other, file)
		}
	}

	return state.GitResult{
		Status:   state.StatusDone,
		Modified: modified,
		Added:    added,
		Deleted:  deleted,
		Other:    other,
	}
}

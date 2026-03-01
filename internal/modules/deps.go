package modules

import (
	"context"
	"strings"
	"time"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

// RunDeps executes `go list -m all` and parses the module list.
// Blocking — call from a tea.Cmd.
func RunDeps(projectDir string) state.DepsResult {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "go", "list", "-m", "all")

	return parseDepsOutput(res)
}

// parseDepsOutput extracts dependency names from `go list -m all` output.
func parseDepsOutput(res services.CommandResult) state.DepsResult {
	if res.Err != nil {
		return state.DepsResult{
			Status: state.StatusError,
			Err:    res.Err.Error(),
		}
	}

	var deps []string
	for _, line := range strings.Split(res.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		deps = append(deps, line)
	}

	// First line is the module itself — skip it if present
	if len(deps) > 1 {
		deps = deps[1:]
	}

	return state.DepsResult{
		Status: state.StatusDone,
		Deps:   deps,
	}
}

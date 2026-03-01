package modules

import (
	"context"
	"strings"
	"time"

	"github.com/chmenegatti/devdash/internal/services"
	"github.com/chmenegatti/devdash/internal/state"
)

// RunLint executes `golangci-lint run` in the project directory and parses
// the output into a LintResult. Blocking — call from a tea.Cmd.
func RunLint(projectDir string) state.LintResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "golangci-lint", "run", "./...")

	return parseLintOutput(res)
}

// parseLintOutput converts golangci-lint output into a LintResult.
func parseLintOutput(res services.CommandResult) state.LintResult {
	combined := res.Stdout + "\n" + res.Stderr

	lines := strings.Split(strings.TrimSpace(combined), "\n")

	var issues []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// golangci-lint prefixes issues with file path; skip meta lines
		if isLintMetaLine(line) {
			continue
		}
		issues = append(issues, line)
	}

	// If the command errored but we found no structured issues,
	// report the raw error.
	if res.Err != nil && len(issues) == 0 {
		return state.LintResult{
			Status: state.StatusError,
			Err:    res.Err.Error(),
		}
	}

	return state.LintResult{
		Status: state.StatusDone,
		Issues: issues,
	}
}

// isLintMetaLine returns true for lines that are not actual lint issues.
func isLintMetaLine(line string) bool {
	lower := strings.ToLower(line)
	for _, prefix := range []string{
		"level=", "msg=", "golangci-lint",
		"run --help", "config file",
	} {
		if strings.Contains(lower, prefix) {
			return true
		}
	}
	return false
}

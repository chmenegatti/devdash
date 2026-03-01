package modules

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

// reCoverageOK matches tested package lines, e.g.:
// ok  github.com/foo/bar 0.123s coverage: 87.3% of statements
var reCoverageOK = regexp.MustCompile(`^ok\s+\S+.*coverage:\s+([\d.]+)%`)

// RunCoverage executes `go test -cover ./...` and parses the aggregate
// coverage percentage. Blocking — call from a tea.Cmd.
func RunCoverage(projectDir string) state.CoverageResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "go", "test", "-cover", "./...")

	return parseCoverageOutput(res)
}

// parseCoverageOutput extracts coverage percentages from go test -cover output.
func parseCoverageOutput(res services.CommandResult) state.CoverageResult {
	combined := res.Stdout + "\n" + res.Stderr

	if res.Err != nil && !strings.Contains(combined, "coverage:") {
		return state.CoverageResult{
			Status: state.StatusError,
			Output: combined,
			Err:    res.Err.Error(),
		}
	}

	lines := strings.Split(combined, "\n")

	var total float64
	var count int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if m := reCoverageOK.FindStringSubmatch(line); m != nil {
			if pct, err := strconv.ParseFloat(m[1], 64); err == nil {
				total += pct
				count++
			}
		}
	}

	var avg float64
	if count > 0 {
		avg = total / float64(count)
	}

	return state.CoverageResult{
		Status:     state.StatusDone,
		Percentage: avg,
		Output:     combined,
	}
}

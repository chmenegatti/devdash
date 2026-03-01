// Package modules provides executable commands for each dashboard feature.
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

// Regular expressions for parsing go test output.
var (
	// Matches: ok  github.com/foo/bar0.123s
	rePkgOK = regexp.MustCompile(`^ok\s+\S+\s+([\d.]+)s`)
	// Matches: FAILgithub.com/foo/bar0.456s
	rePkgFAIL = regexp.MustCompile(`^FAIL\s+\S+\s+([\d.]+)s`)
	// Matches: ?   github.com/foo/bar[no test files]
	rePkgSkip = regexp.MustCompile(`^\?\s+\S+`)
)

// RunTests executes `go test ./...` in the project directory and returns a
// fully populated TestsResult. This function is blocking and should be called
// from a Bubble Tea Cmd (goroutine).
func RunTests(projectDir string) state.TestsResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "go", "test", "-v", "-count=1", "./...")

	return parseTestOutput(res)
}

// parseTestOutput converts command output into a TestsResult.
func parseTestOutput(res services.CommandResult) state.TestsResult {
	// Combine stdout and stderr — go test writes some output to stderr.
	combined := res.Stdout + "\n" + res.Stderr

	lines := strings.Split(combined, "\n")

	var (
		passed    = true
		pkgCount  int
		totalSecs float64
		hasFail   bool
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if m := rePkgOK.FindStringSubmatch(line); m != nil {
			pkgCount++
			if dur, err := strconv.ParseFloat(m[1], 64); err == nil {
				totalSecs += dur
			}
			continue
		}

		if m := rePkgFAIL.FindStringSubmatch(line); m != nil {
			pkgCount++
			hasFail = true
			if dur, err := strconv.ParseFloat(m[1], 64); err == nil {
				totalSecs += dur
			}
			continue
		}

		if rePkgSkip.MatchString(line) {
			// Package with [no test files] — count but don't affect pass/fail
			pkgCount++
			continue
		}
	}

	if hasFail || res.Err != nil {
		passed = false
	}

	// If we didn't find any package lines at all, still report error context.
	if pkgCount == 0 && res.Err != nil {
		return state.TestsResult{
			Status: state.StatusError,
			Passed: false,
			Output: combined,
			Err:    res.Err.Error(),
		}
	}

	duration := time.Duration(totalSecs * float64(time.Second))

	return state.TestsResult{
		Status:   state.StatusDone,
		Passed:   passed,
		Packages: pkgCount,
		Duration: duration,
		Output:   combined,
	}
}

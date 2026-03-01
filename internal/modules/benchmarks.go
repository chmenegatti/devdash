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

// reBench matches: BenchmarkRetry-8    1200000    1200 ns/op
var reBench = regexp.MustCompile(`^(Benchmark\S+)\s+(\d+)\s+([\d.]+)\s+ns/op`)

// RunBenchmarks executes `go test -bench=. -benchmem ./...` and parses results.
// Blocking — call from a tea.Cmd.
func RunBenchmarks(projectDir string) state.BenchmarkResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	res := services.RunCommand(ctx, projectDir, "go", "test", "-bench=.", "-benchmem", "-run=^$", "./...")

	return parseBenchmarkOutput(res)
}

// parseBenchmarkOutput extracts benchmark entries from go test -bench output.
func parseBenchmarkOutput(res services.CommandResult) state.BenchmarkResult {
	combined := res.Stdout + "\n" + res.Stderr

	// go test -bench can fail partly — still parse what we got
	lines := strings.Split(combined, "\n")

	var entries []state.BenchmarkEntry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if m := reBench.FindStringSubmatch(line); m != nil {
			iters, _ := strconv.Atoi(m[2])
			nsPerOp, _ := strconv.ParseFloat(m[3], 64)
			entries = append(entries, state.BenchmarkEntry{
				Name:       m[1],
				Iterations: iters,
				NsPerOp:    nsPerOp,
			})
		}
	}

	if len(entries) == 0 && res.Err != nil {
		return state.BenchmarkResult{
			Status: state.StatusError,
			Output: combined,
			Err:    res.Err.Error(),
		}
	}

	return state.BenchmarkResult{
		Status:  state.StatusDone,
		Entries: entries,
		Output:  combined,
	}
}

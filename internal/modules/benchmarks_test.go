package modules

import (
	"fmt"
	"testing"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

func TestParseBenchmarkOutput_Multiple(t *testing.T) {
	res := services.CommandResult{
		Stdout: "goos: linux\ngoarch: amd64\nBenchmarkRetry-8          1200000              1200 ns/op\nBenchmarkParse-8          5000000               300.5 ns/op\nPASS\nok  \tgithub.com/foo/bar\t3.200s\n",
	}
	got := parseBenchmarkOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	if len(got.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got.Entries))
	}
	if got.Entries[0].Name != "BenchmarkRetry-8" {
		t.Errorf("expected BenchmarkRetry-8, got %s", got.Entries[0].Name)
	}
	if got.Entries[0].Iterations != 1200000 {
		t.Errorf("expected 1200000 iters, got %d", got.Entries[0].Iterations)
	}
	if got.Entries[0].NsPerOp != 1200 {
		t.Errorf("expected 1200 ns/op, got %.1f", got.Entries[0].NsPerOp)
	}
	if got.Entries[1].NsPerOp != 300.5 {
		t.Errorf("expected 300.5 ns/op, got %.1f", got.Entries[1].NsPerOp)
	}
}

func TestParseBenchmarkOutput_NoBenchmarks(t *testing.T) {
	res := services.CommandResult{
		Stdout: "?   \tgithub.com/foo/bar\t[no test files]\n",
	}
	got := parseBenchmarkOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	if len(got.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got.Entries))
	}
}

func TestParseBenchmarkOutput_Error(t *testing.T) {
	res := services.CommandResult{
		Stderr: "build failed",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseBenchmarkOutput(res)

	if got.Status != state.StatusError {
		t.Errorf("expected StatusError, got %v", got.Status)
	}
}

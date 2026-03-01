package modules

import (
	"fmt"
	"testing"
	"time"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

func TestParseTestOutput_AllPass(t *testing.T) {
	res := services.CommandResult{
		Stdout: "ok  \tgithub.com/foo/bar\t0.250s\nok  \tgithub.com/foo/baz\t1.100s\n",
	}
	got := parseTestOutput(res)

	if got.Status != state.StatusDone {
		t.Errorf("expected StatusDone, got %v", got.Status)
	}
	if !got.Passed {
		t.Error("expected Passed=true")
	}
	if got.Packages != 2 {
		t.Errorf("expected 2 packages, got %d", got.Packages)
	}
	// 0.250 + 1.100 = 1.350s
	want := time.Duration(1350 * time.Millisecond)
	if got.Duration != want {
		t.Errorf("expected duration %v, got %v", want, got.Duration)
	}
}

func TestParseTestOutput_WithFailure(t *testing.T) {
	res := services.CommandResult{
		Stdout: "ok  \tgithub.com/foo/bar\t0.100s\nFAIL\tgithub.com/foo/baz\t0.200s\n",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseTestOutput(res)

	if got.Passed {
		t.Error("expected Passed=false")
	}
	if got.Packages != 2 {
		t.Errorf("expected 2 packages, got %d", got.Packages)
	}
}

func TestParseTestOutput_NoTestFiles(t *testing.T) {
	res := services.CommandResult{
		Stdout: "?   \tgithub.com/foo/bar\t[no test files]\n",
	}
	got := parseTestOutput(res)

	if got.Packages != 1 {
		t.Errorf("expected 1 package, got %d", got.Packages)
	}
	if !got.Passed {
		t.Error("expected Passed=true for skip-only")
	}
}

func TestParseTestOutput_Error(t *testing.T) {
	res := services.CommandResult{
		Stderr: "cannot find module",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseTestOutput(res)

	if got.Status != state.StatusError {
		t.Errorf("expected StatusError, got %v", got.Status)
	}
}

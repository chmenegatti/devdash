package modules

import (
	"testing"

	"github.com/chmenegatti/devdash/internal/state"
)

func TestRunTests_InvalidProjectDir(t *testing.T) {
	got := RunTests("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

func TestRunCoverage_InvalidProjectDir(t *testing.T) {
	got := RunCoverage("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

func TestRunLint_InvalidProjectDir(t *testing.T) {
	got := RunLint("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

func TestRunBenchmarks_InvalidProjectDir(t *testing.T) {
	got := RunBenchmarks("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

func TestRunDeps_InvalidProjectDir(t *testing.T) {
	got := RunDeps("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

func TestRunGitStatus_InvalidProjectDir(t *testing.T) {
	got := RunGitStatus("/tmp/devdash-dir-does-not-exist")
	if got.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", got.Status)
	}
}

package modules

import (
	"testing"

	"github.com/cesar/devdash/internal/state"
)

func TestRunBinarySize_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	result := RunBinarySize("../../")
	if result.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v (err: %s)", result.Status, result.Err)
	}
	if result.Size <= 0 {
		t.Fatalf("expected positive binary size, got %d", result.Size)
	}
	t.Logf("binary size: %d bytes", result.Size)
}

func TestRunBinarySize_BuildError(t *testing.T) {
	result := RunBinarySize("/tmp/nonexistent_devdash_proj")
	if result.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", result.Status)
	}
	if result.Err == "" {
		t.Fatal("expected error message, got empty")
	}
}

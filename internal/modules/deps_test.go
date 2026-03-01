package modules

import (
	"fmt"
	"testing"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

func TestParseDepsOutput_MultipleModules(t *testing.T) {
	res := services.CommandResult{
		Stdout: "github.com/cesar/devdash\ngithub.com/charmbracelet/bubbletea v1.3.10\ngithub.com/charmbracelet/lipgloss v1.1.0\n",
	}
	result := parseDepsOutput(res)
	if result.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", result.Status)
	}
	// First line (self module) should be skipped
	if len(result.Deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(result.Deps), result.Deps)
	}
}

func TestParseDepsOutput_NoExternalDeps(t *testing.T) {
	res := services.CommandResult{
		Stdout: "github.com/cesar/devdash\n",
	}
	result := parseDepsOutput(res)
	if result.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", result.Status)
	}
	// Only module itself — after skipping, should be just the one line
	if len(result.Deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(result.Deps))
	}
}

func TestParseDepsOutput_Error(t *testing.T) {
	res := services.CommandResult{
		Err: fmt.Errorf("module lookup failed"),
	}
	result := parseDepsOutput(res)
	if result.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", result.Status)
	}
}

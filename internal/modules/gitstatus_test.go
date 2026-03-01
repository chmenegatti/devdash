package modules

import (
	"fmt"
	"testing"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

func TestParseGitOutput_Mixed(t *testing.T) {
	res := services.CommandResult{
		Stdout: " M internal/app/app.go\nA  README.md\n D old.go\n?? scratch.txt\n",
	}
	result := parseGitOutput(res)
	if result.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", result.Status)
	}
	if len(result.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(result.Modified))
	}
	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added))
	}
	if len(result.Deleted) != 1 {
		t.Fatalf("expected 1 deleted, got %d", len(result.Deleted))
	}
	if len(result.Other) != 1 {
		t.Fatalf("expected 1 other, got %d", len(result.Other))
	}
}

func TestParseGitOutput_Clean(t *testing.T) {
	res := services.CommandResult{
		Stdout: "",
	}
	result := parseGitOutput(res)
	if result.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", result.Status)
	}
	total := len(result.Modified) + len(result.Added) + len(result.Deleted) + len(result.Other)
	if total != 0 {
		t.Fatalf("expected 0 changes, got %d", total)
	}
}

func TestParseGitOutput_Error(t *testing.T) {
	res := services.CommandResult{
		Err: fmt.Errorf("not a git repository"),
	}
	result := parseGitOutput(res)
	if result.Status != state.StatusError {
		t.Fatalf("expected StatusError, got %v", result.Status)
	}
}

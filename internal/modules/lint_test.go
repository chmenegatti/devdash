package modules

import (
	"fmt"
	"testing"

	"github.com/chmenegatti/devdash/internal/services"
	"github.com/chmenegatti/devdash/internal/state"
)

func TestParseLintOutput_Clean(t *testing.T) {
	res := services.CommandResult{
		Stdout: "",
		Stderr: "",
	}
	got := parseLintOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	if len(got.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(got.Issues))
	}
}

func TestParseLintOutput_WithIssues(t *testing.T) {
	res := services.CommandResult{
		Stdout: "main.go:10:2: unused variable 'x' (deadcode)\nmain.go:15:4: error return value not checked (errcheck)\n",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseLintOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	if len(got.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(got.Issues), got.Issues)
	}
}

func TestParseLintOutput_FiltersMeta(t *testing.T) {
	res := services.CommandResult{
		Stderr: "level=warning msg=\"some config warning\"\nmain.go:5:1: exported function Foo should have comment (golint)\n",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseLintOutput(res)

	if len(got.Issues) != 1 {
		t.Errorf("expected 1 issue after filtering meta, got %d: %v", len(got.Issues), got.Issues)
	}
}

func TestParseLintOutput_CommandNotFound(t *testing.T) {
	res := services.CommandResult{
		Err: fmt.Errorf("exec: \"golangci-lint\": executable file not found in $PATH"),
	}
	got := parseLintOutput(res)

	if got.Status != state.StatusError {
		t.Errorf("expected StatusError, got %v", got.Status)
	}
}

package modules

import (
	"fmt"
	"testing"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
)

func TestParseCoverageOutput_MultiPackage(t *testing.T) {
	res := services.CommandResult{
		Stdout: "github.com/foo/cmd\t\tcoverage: 0.0% of statements\n" +
			"ok  \tgithub.com/foo/bar\t0.1s\tcoverage: 80.0% of statements\n" +
			"ok  \tgithub.com/foo/baz\t0.2s\tcoverage: 90.0% of statements\n",
	}
	got := parseCoverageOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	// average of 80 + 90 = 85
	if got.Percentage != 85.0 {
		t.Errorf("expected 85.0%%, got %.1f%%", got.Percentage)
	}
}

func TestParseCoverageOutput_SinglePackage(t *testing.T) {
	res := services.CommandResult{
		Stdout: "ok  \tgithub.com/foo/bar\t0.1s\tcoverage: 42.5% of statements\n",
	}
	got := parseCoverageOutput(res)

	if got.Percentage != 42.5 {
		t.Errorf("expected 42.5%%, got %.1f%%", got.Percentage)
	}
}

func TestParseCoverageOutput_NoCoverage(t *testing.T) {
	res := services.CommandResult{
		Stdout: "?   \tgithub.com/foo/bar\t[no test files]\n",
	}
	got := parseCoverageOutput(res)

	if got.Status != state.StatusDone {
		t.Fatalf("expected StatusDone, got %v", got.Status)
	}
	if got.Percentage != 0 {
		t.Errorf("expected 0%%, got %.1f%%", got.Percentage)
	}
}

func TestParseCoverageOutput_Error(t *testing.T) {
	res := services.CommandResult{
		Stderr: "cannot find module",
		Err:    fmt.Errorf("exit status 1"),
	}
	got := parseCoverageOutput(res)

	if got.Status != state.StatusError {
		t.Errorf("expected StatusError, got %v", got.Status)
	}
}

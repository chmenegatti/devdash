package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cesar/devdash/internal/state"
)

func TestGenerateMarkdownReport_IncludesAllSections(t *testing.T) {
	ds := state.Dashboard{
		ProjectName: "devdash",
		ProjectDir:  "/tmp/devdash",
		Tests: state.TestsResult{
			Status:   state.StatusDone,
			Passed:   true,
			Packages: 7,
			Duration: 2 * time.Second,
			Output:   "ok github.com/cesar/devdash 2.0s",
		},
		Coverage: state.CoverageResult{
			Status:     state.StatusDone,
			Percentage: 97.1,
		},
		Lint: state.LintResult{
			Status: state.StatusDone,
			Issues: []string{"internal/x.go:10: something"},
		},
		Benchmarks: state.BenchmarkResult{
			Status: state.StatusDone,
			Entries: []state.BenchmarkEntry{{
				Name:       "BenchmarkParse-8",
				Iterations: 1000,
				NsPerOp:    321.4,
			}},
		},
		Binary: state.BinaryResult{Status: state.StatusDone, Size: 4200000},
		Deps:   state.DepsResult{Status: state.StatusDone, Deps: []string{"github.com/charmbracelet/lipgloss v1.1.0"}},
		Git:    state.GitResult{Status: state.StatusDone, Modified: []string{"internal/app/app.go"}},
	}

	report := generateMarkdownReport(ds, time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC))

	mustContain := []string{
		"# 🧾 devdash Report",
		"## 📌 Executive Summary",
		"## 🧪 Tests",
		"## 📊 Coverage",
		"## 🔍 Lint",
		"## ⚡ Benchmarks",
		"## 📦 Binary",
		"## 🌿 Git Status",
		"## 📚 Dependencies",
		"BenchmarkParse-8",
		"97.1%",
	}

	for _, token := range mustContain {
		if !strings.Contains(report, token) {
			t.Fatalf("expected report to contain %q", token)
		}
	}
}

func TestGenerateReportFile_WritesMarkdownFile(t *testing.T) {
	tmp := t.TempDir()
	ds := state.Dashboard{ProjectName: "demo", ProjectDir: tmp}

	path, err := GenerateReportFile(tmp, ds)
	if err != nil {
		t.Fatalf("GenerateReportFile returned error: %v", err)
	}

	if filepath.Ext(path) != ".md" {
		t.Fatalf("expected .md file, got %s", path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read report file: %v", err)
	}
	if !strings.Contains(string(b), "devdash Report") {
		t.Fatalf("report content is missing expected title")
	}
}

func TestExtractTestCaseResults(t *testing.T) {
	output := strings.Join([]string{
		"=== RUN   TestAlpha",
		"--- PASS: TestAlpha (0.00s)",
		"=== RUN   TestBeta",
		"--- FAIL: TestBeta (0.00s)",
		"FAIL",
	}, "\n")

	passed, failed := extractTestCaseResults(output)
	if len(passed) != 1 || passed[0] != "TestAlpha" {
		t.Fatalf("expected [TestAlpha], got %v", passed)
	}
	if len(failed) != 1 || failed[0] != "TestBeta" {
		t.Fatalf("expected [TestBeta], got %v", failed)
	}
}

func TestGenerateMarkdownReport_IncludesTestCaseIcons(t *testing.T) {
	ds := state.Dashboard{
		ProjectName: "devdash",
		ProjectDir:  "/tmp/devdash",
		Tests: state.TestsResult{
			Status: state.StatusDone,
			Passed: false,
			Output: strings.Join([]string{
				"--- PASS: TestOne (0.00s)",
				"--- FAIL: TestTwo (0.01s)",
			}, "\n"),
		},
	}

	report := generateMarkdownReport(ds, time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC))
	if !strings.Contains(report, "- ✅ TestOne") {
		t.Fatalf("expected report to include pass icon line")
	}
	if !strings.Contains(report, "- ❌ TestTwo") {
		t.Fatalf("expected report to include fail icon line")
	}
}

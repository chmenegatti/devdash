package modules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cesar/devdash/internal/services"
)

func TestParserGoldens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		stdoutFile   string
		stderrFile   string
		errText      string
		run          func(services.CommandResult) string
		expectedFile string
	}{
		{
			name:         "tests_all_pass",
			stdoutFile:   "tests_all_pass.stdout",
			run:          snapshotTestsResult,
			expectedFile: "tests_all_pass.expected",
		},
		{
			name:         "tests_with_failure",
			stdoutFile:   "tests_with_failure.stdout",
			errText:      "exit status 1",
			run:          snapshotTestsResult,
			expectedFile: "tests_with_failure.expected",
		},
		{
			name:         "coverage_multi",
			stdoutFile:   "coverage_multi.stdout",
			run:          snapshotCoverageResult,
			expectedFile: "coverage_multi.expected",
		},
		{
			name:         "lint_with_meta",
			stdoutFile:   "lint_with_meta.stdout",
			stderrFile:   "lint_with_meta.stderr",
			errText:      "exit status 1",
			run:          snapshotLintResult,
			expectedFile: "lint_with_meta.expected",
		},
		{
			name:         "bench_two_entries",
			stdoutFile:   "bench_two_entries.stdout",
			run:          snapshotBenchResult,
			expectedFile: "bench_two_entries.expected",
		},
		{
			name:         "deps_two_modules",
			stdoutFile:   "deps_two_modules.stdout",
			run:          snapshotDepsResult,
			expectedFile: "deps_two_modules.expected",
		},
		{
			name:         "git_mixed",
			stdoutFile:   "git_mixed.stdout",
			run:          snapshotGitResult,
			expectedFile: "git_mixed.expected",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := services.CommandResult{
				Stdout: mustReadFixture(t, tc.stdoutFile),
				Stderr: mustReadFixtureIfExists(t, tc.stderrFile),
			}
			if tc.errText != "" {
				res.Err = errors.New(tc.errText)
			}

			got := strings.TrimSpace(tc.run(res))
			want := strings.TrimSpace(mustReadFixture(t, tc.expectedFile))
			if got != want {
				t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", tc.name, got, want)
			}
		})
	}
}

func mustReadFixture(t *testing.T, fileName string) string {
	t.Helper()
	if strings.TrimSpace(fileName) == "" {
		return ""
	}
	path := filepath.Join("testdata", "parsers", fileName)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	return string(b)
}

func mustReadFixtureIfExists(t *testing.T, fileName string) string {
	t.Helper()
	if strings.TrimSpace(fileName) == "" {
		return ""
	}
	path := filepath.Join("testdata", "parsers", fileName)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	return string(b)
}

func snapshotTestsResult(res services.CommandResult) string {
	r := parseTestOutput(res)
	return strings.Join([]string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("passed=%t", r.Passed),
		fmt.Sprintf("testCases=%d", r.TestCases),
		fmt.Sprintf("failedTests=%d", r.FailedTests),
		fmt.Sprintf("packages=%d", r.Packages),
		fmt.Sprintf("duration=%s", r.Duration),
	}, "\n")
}

func snapshotCoverageResult(res services.CommandResult) string {
	r := parseCoverageOutput(res)
	return strings.Join([]string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("percentage=%.1f", r.Percentage),
	}, "\n")
}

func snapshotLintResult(res services.CommandResult) string {
	r := parseLintOutput(res)
	lines := []string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("issues=%d", len(r.Issues)),
	}
	for _, issue := range r.Issues {
		lines = append(lines, "- "+issue)
	}
	return strings.Join(lines, "\n")
}

func snapshotBenchResult(res services.CommandResult) string {
	r := parseBenchmarkOutput(res)
	lines := []string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("entries=%d", len(r.Entries)),
	}
	for _, e := range r.Entries {
		lines = append(lines, fmt.Sprintf("- %s|%d|%.1f", e.Name, e.Iterations, e.NsPerOp))
	}
	return strings.Join(lines, "\n")
}

func snapshotDepsResult(res services.CommandResult) string {
	r := parseDepsOutput(res)
	lines := []string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("deps=%d", len(r.Deps)),
	}
	for _, dep := range r.Deps {
		lines = append(lines, "- "+dep)
	}
	return strings.Join(lines, "\n")
}

func snapshotGitResult(res services.CommandResult) string {
	r := parseGitOutput(res)
	return strings.Join([]string{
		fmt.Sprintf("status=%s", r.Status.String()),
		fmt.Sprintf("modified=%d", len(r.Modified)),
		fmt.Sprintf("added=%d", len(r.Added)),
		fmt.Sprintf("deleted=%d", len(r.Deleted)),
		fmt.Sprintf("other=%d", len(r.Other)),
	}, "\n")
}

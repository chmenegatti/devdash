// Package state provides centralized state management for the dashboard.
package state

import "time"

// Status represents the state of an operation.
type Status int

const (
	StatusIdle    Status = iota // Not yet run
	StatusRunning               // Currently executing
	StatusDone                  // Completed successfully
	StatusError                 // Completed with error
)

// String returns a human-readable label for the status.
func (s Status) String() string {
	switch s {
	case StatusRunning:
		return "Running…"
	case StatusDone:
		return "Done"
	case StatusError:
		return "Error"
	default:
		return "Idle"
	}
}

// TestsResult holds the outcome of running go test.
type TestsResult struct {
	Status   Status
	Passed   bool
	Packages int
	Duration time.Duration
	Output   string
	Err      string
}

// CoverageResult holds coverage information.
type CoverageResult struct {
	Status     Status
	Percentage float64
	Output     string
	Err        string
}

// LintResult holds linting information.
type LintResult struct {
	Status Status
	Issues []string
	Err    string
}

// BenchmarkEntry represents a single benchmark line.
type BenchmarkEntry struct {
	Name       string
	Iterations int
	NsPerOp    float64
}

// BenchmarkResult holds benchmark information.
type BenchmarkResult struct {
	Status  Status
	Entries []BenchmarkEntry
	Output  string
	Err     string
}

// BinaryResult holds binary size information.
type BinaryResult struct {
	Status Status
	Size   int64 // bytes
	Err    string
}

// DepsResult holds dependency information.
type DepsResult struct {
	Status Status
	Deps   []string
	Err    string
}

// GitResult holds git status information.
type GitResult struct {
	Status   Status
	Modified []string
	Added    []string
	Deleted  []string
	Other    []string
	Err      string
}

// Dashboard is the central state container for the entire application.
type Dashboard struct {
	ProjectDir  string
	ProjectName string

	Tests      TestsResult
	Coverage   CoverageResult
	Lint       LintResult
	Benchmarks BenchmarkResult
	Binary     BinaryResult
	Deps       DepsResult
	Git        GitResult
}

// New creates a new Dashboard state for the given project directory.
func New(projectDir, projectName string) *Dashboard {
	return &Dashboard{
		ProjectDir:  projectDir,
		ProjectName: projectName,
	}
}

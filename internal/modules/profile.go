package modules

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chmenegatti/devdash/internal/services"
	"github.com/chmenegatti/devdash/internal/state"
	"github.com/google/pprof/profile"
)

const (
	flamegraphMaxNodes = 120
	flamegraphBarWidth = 22
)

const noSamplesFlamegraph = "No CPU samples captured.\n\nTry one of the following:\n- run profiling in a package with heavier tests\n- run profiling while benchmarks exist in the package\n- run profiling multiple times (sampling can miss very short runs)"

type flameNode struct {
	name     string
	value    int64
	children map[string]*flameNode
}

func newFlameNode(name string) *flameNode {
	return &flameNode{
		name:     name,
		children: make(map[string]*flameNode),
	}
}

// RunCPUProfile executes go test with cpuprofile on one package and generates
// an inline Unicode flamegraph for terminal rendering.
func RunCPUProfile(projectDir string) state.ProfileResult {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	targetPkgs, err := detectProfileTargetPackages(ctx, projectDir)
	if err != nil {
		return state.ProfileResult{
			Status: state.StatusError,
			Err:    err.Error(),
		}
	}

	var lastErr error
	for _, targetPkg := range targetPkgs {
		flame, total, unit, err := runProfileForPackage(ctx, projectDir, targetPkg)
		if err != nil {
			lastErr = err
			continue
		}
		if total == 0 {
			continue
		}

		return state.ProfileResult{
			Status:        state.StatusDone,
			TargetPackage: targetPkg,
			TotalSamples:  total,
			SampleUnit:    unit,
			Flamegraph:    flame,
		}
	}

	if lastErr != nil {
		return state.ProfileResult{
			Status: state.StatusError,
			Err:    lastErr.Error(),
		}
	}

	return state.ProfileResult{
		Status:        state.StatusDone,
		TargetPackage: targetPkgs[0],
		TotalSamples:  0,
		SampleUnit:    "samples",
		Flamegraph:    noSamplesFlamegraph,
	}
}

func runProfileForPackage(ctx context.Context, projectDir, targetPkg string) (flamegraph string, total int64, sampleUnit string, err error) {
	tmp, err := os.CreateTemp("", "devdash-cpu-*.pprof")
	if err != nil {
		return "", 0, "", fmt.Errorf("create temp profile: %w", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	res := services.RunCommand(ctx, projectDir, "go", "test", "-count=1", "-cpuprofile="+tmpPath, targetPkg)
	if res.Err != nil {
		return "", 0, "", fmt.Errorf("profile %s: %w\n%s", targetPkg, res.Err, strings.TrimSpace(res.Stderr))
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return "", 0, "", fmt.Errorf("open cpu profile: %w", err)
	}
	defer func() { _ = f.Close() }()

	prof, err := profile.Parse(f)
	if err != nil {
		return "", 0, "", fmt.Errorf("parse cpu profile: %w", err)
	}

	flame, total := buildInlineFlamegraph(prof)
	return flame, total, profileSampleUnit(prof), nil
}

func detectProfileTargetPackages(ctx context.Context, projectDir string) ([]string, error) {
	tmpl := `{{if or (gt (len .TestGoFiles) 0) (gt (len .XTestGoFiles) 0)}}{{len .TestGoFiles}} {{len .XTestGoFiles}} {{.ImportPath}}{{end}}`
	list := services.RunCommand(ctx, projectDir, "go", "list", "-f", tmpl, "./...")
	if list.Err != nil {
		return nil, fmt.Errorf("detect profile package: %w\n%s", list.Err, strings.TrimSpace(list.Stderr))
	}

	type candidate struct {
		pkg   string
		score int
	}

	var candidates []candidate

	for _, line := range strings.Split(list.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		testsN, err1 := strconv.Atoi(parts[0])
		xTestsN, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}
		pkg := strings.Join(parts[2:], " ")
		if pkg == "" {
			continue
		}
		candidates = append(candidates, candidate{pkg: pkg, score: testsN + xTestsN})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no Go package with tests found for CPU profiling")
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].pkg < candidates[j].pkg
		}
		return candidates[i].score > candidates[j].score
	})

	pkgs := make([]string, 0, len(candidates))
	for _, c := range candidates {
		pkgs = append(pkgs, c.pkg)
	}

	return pkgs, nil
}

func buildInlineFlamegraph(prof *profile.Profile) (string, int64) {
	root := newFlameNode("root")
	var total int64

	for _, sample := range prof.Sample {
		if len(sample.Value) == 0 {
			continue
		}
		val := sample.Value[0]
		if val <= 0 {
			continue
		}
		total += val

		stack := sampleFunctionStack(sample)
		node := root
		node.value += val
		for _, fn := range stack {
			child := node.children[fn]
			if child == nil {
				child = newFlameNode(fn)
				node.children[fn] = child
			}
			child.value += val
			node = child
		}
	}

	if total == 0 {
		return "", 0
	}

	return renderFlamegraphTree(root, total, flamegraphMaxNodes), total
}

func sampleFunctionStack(sample *profile.Sample) []string {
	stack := make([]string, 0, len(sample.Location))
	for i := len(sample.Location) - 1; i >= 0; i-- {
		name := locationFunctionName(sample.Location[i])
		if name == "" {
			continue
		}
		stack = append(stack, shortFunctionName(name))
	}
	return stack
}

func locationFunctionName(loc *profile.Location) string {
	if loc == nil {
		return ""
	}
	for _, line := range loc.Line {
		if line.Function != nil && line.Function.Name != "" {
			return line.Function.Name
		}
	}
	return ""
}

func shortFunctionName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 && idx < len(name)-1 {
		name = name[idx+1:]
	}
	return name
}

func renderFlamegraphTree(root *flameNode, total int64, maxNodes int) string {
	if root == nil || total <= 0 {
		return ""
	}

	lines := []string{"CPU Flamegraph (inline)", ""}
	count := 0

	var walk func(node *flameNode, depth int)
	walk = func(node *flameNode, depth int) {
		if node == nil || count >= maxNodes {
			return
		}

		children := sortedChildren(node)
		for _, child := range children {
			if count >= maxNodes {
				return
			}
			pct := (float64(child.value) / float64(total)) * 100
			indent := strings.Repeat("  ", depth)
			bar := unicodeBar(child.value, total, flamegraphBarWidth)
			lines = append(lines, fmt.Sprintf("%s%s %-38s %6.2f%%", indent, bar, truncateText(child.name, 38), pct))
			count++
			walk(child, depth+1)
		}
	}

	walk(root, 0)
	if count == 0 {
		return "(no flamegraph nodes)"
	}

	if count >= maxNodes {
		lines = append(lines, "", "… truncated …")
	}

	return strings.Join(lines, "\n")
}

func sortedChildren(node *flameNode) []*flameNode {
	items := make([]*flameNode, 0, len(node.children))
	for _, child := range node.children {
		items = append(items, child)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].value == items[j].value {
			return items[i].name < items[j].name
		}
		return items[i].value > items[j].value
	})
	return items
}

func unicodeBar(value, total int64, width int) string {
	if width < 1 {
		width = 1
	}
	if total <= 0 || value <= 0 {
		return strings.Repeat("░", width)
	}
	filled := int((float64(value) / float64(total)) * float64(width))
	if filled < 1 {
		filled = 1
	}
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func profileSampleUnit(prof *profile.Profile) string {
	if prof == nil || len(prof.SampleType) == 0 {
		return "samples"
	}
	st := prof.SampleType[0]
	if st.Unit != "" {
		return st.Unit
	}
	if st.Type != "" {
		return st.Type
	}
	return "samples"
}

func truncateText(s string, max int) string {
	if max < 4 {
		max = 4
	}
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

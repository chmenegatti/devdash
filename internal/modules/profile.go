package modules

import (
	"context"
	"fmt"
	"os"
	"sort"
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

	targetPkg, err := detectProfileTargetPackage(ctx, projectDir)
	if err != nil {
		return state.ProfileResult{
			Status: state.StatusError,
			Err:    err.Error(),
		}
	}

	tmp, err := os.CreateTemp("", "devdash-cpu-*.pprof")
	if err != nil {
		return state.ProfileResult{
			Status: state.StatusError,
			Err:    "create temp profile: " + err.Error(),
		}
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	res := services.RunCommand(ctx, projectDir, "go", "test", "-count=1", "-cpuprofile="+tmpPath, targetPkg)
	if res.Err != nil {
		return state.ProfileResult{
			Status:        state.StatusError,
			TargetPackage: targetPkg,
			Err:           strings.TrimSpace(res.Err.Error() + "\n" + res.Stderr),
		}
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return state.ProfileResult{
			Status:        state.StatusError,
			TargetPackage: targetPkg,
			Err:           "open cpu profile: " + err.Error(),
		}
	}
	defer func() { _ = f.Close() }()

	prof, err := profile.Parse(f)
	if err != nil {
		return state.ProfileResult{
			Status:        state.StatusError,
			TargetPackage: targetPkg,
			Err:           "parse cpu profile: " + err.Error(),
		}
	}

	flame, total := buildInlineFlamegraph(prof)
	if total == 0 {
		return state.ProfileResult{
			Status:        state.StatusError,
			TargetPackage: targetPkg,
			Err:           "profile has no CPU samples",
		}
	}

	return state.ProfileResult{
		Status:        state.StatusDone,
		TargetPackage: targetPkg,
		TotalSamples:  total,
		SampleUnit:    profileSampleUnit(prof),
		Flamegraph:    flame,
	}
}

func detectProfileTargetPackage(ctx context.Context, projectDir string) (string, error) {
	tmpl := `{{if or (gt (len .TestGoFiles) 0) (gt (len .XTestGoFiles) 0)}}{{.ImportPath}}{{end}}`
	list := services.RunCommand(ctx, projectDir, "go", "list", "-f", tmpl, "./...")
	if list.Err != nil {
		return "", fmt.Errorf("detect profile package: %w\n%s", list.Err, strings.TrimSpace(list.Stderr))
	}

	for _, line := range strings.Split(list.Stdout, "\n") {
		pkg := strings.TrimSpace(line)
		if pkg != "" {
			return pkg, nil
		}
	}

	return "", fmt.Errorf("no Go package with tests found for CPU profiling")
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

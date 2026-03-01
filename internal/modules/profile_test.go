package modules

import (
	"strings"
	"testing"
)

func TestUnicodeBar(t *testing.T) {
	bar := unicodeBar(50, 100, 10)
	if strings.Count(bar, "█") == 0 {
		t.Fatalf("expected at least one filled block, got %q", bar)
	}
	if len([]rune(bar)) != 10 {
		t.Fatalf("expected bar width 10, got %d", len([]rune(bar)))
	}
}

func TestRenderFlamegraphTree_SortedByWeight(t *testing.T) {
	root := newFlameNode("root")
	root.children["small"] = &flameNode{name: "small", value: 10, children: map[string]*flameNode{}}
	root.children["big"] = &flameNode{name: "big", value: 90, children: map[string]*flameNode{}}

	out := renderFlamegraphTree(root, 100, 10)
	idxBig := strings.Index(out, "big")
	idxSmall := strings.Index(out, "small")
	if idxBig == -1 || idxSmall == -1 {
		t.Fatalf("expected both nodes in output, got:\n%s", out)
	}
	if idxBig > idxSmall {
		t.Fatalf("expected heavier node first, got:\n%s", out)
	}
}

func TestNoSamplesFlamegraph_HasGuidance(t *testing.T) {
	if !strings.Contains(noSamplesFlamegraph, "No CPU samples captured") {
		t.Fatalf("expected guidance message, got %q", noSamplesFlamegraph)
	}
	if !strings.Contains(noSamplesFlamegraph, "heavier tests") {
		t.Fatalf("expected actionable guidance, got %q", noSamplesFlamegraph)
	}
}

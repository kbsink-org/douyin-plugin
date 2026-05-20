package xhsconv

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

func TestParser_Parse(t *testing.T) {
	html := mustReadTestHTML(t, "xiaohongshu_real.html")
	res, err := NewParser().Parse(context.Background(), &core.FetchResult{
		URL:  "https://www.xiaohongshu.com/explore/69eca7e800000000230072ba",
		HTML: html,
	}, "output")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if strings.TrimSpace(res.Title) == "" {
		t.Fatalf("expected non-empty title from real snapshot")
	}
	if strings.TrimSpace(res.Markdown) == "" {
		t.Fatalf("expected non-empty markdown from real snapshot")
	}
	if len(res.Images) < 2 {
		t.Fatalf("expected at least 2 note images from multi-image snapshot, got %d", len(res.Images))
	}
	if !strings.Contains(res.Markdown, "![](http") && !strings.Contains(res.Markdown, "![](https") {
		t.Fatalf("expected markdown to embed images as ![](...), got prefix %q", truncate(res.Markdown, 200))
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func mustReadTestHTML(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join("testdata", name)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read testdata %q: %v", p, err)
	}
	return string(b)
}

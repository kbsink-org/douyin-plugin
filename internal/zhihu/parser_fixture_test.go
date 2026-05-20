package zhihuconv

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

func TestParser_Fixture(t *testing.T) {
	html := mustReadTestHTML(t, "zhihu_zhuanlan_sample.html")
	res, err := NewParser().Parse(context.Background(), &core.FetchResult{
		URL:  "https://zhuanlan.zhihu.com/p/26954669801",
		HTML: html,
	}, "output")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if strings.TrimSpace(res.Title) == "" {
		t.Fatalf("expected non-empty title")
	}
	if !strings.Contains(res.Markdown, "示例段落") {
		t.Fatalf("expected article body in markdown: %q", res.Markdown)
	}
	if res.AccountName != "专栏作者" {
		t.Fatalf("unexpected account: %q", res.AccountName)
	}
	if len(res.Images) < 1 {
		t.Fatalf("expected at least one image")
	}
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

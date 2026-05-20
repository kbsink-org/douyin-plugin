package bilibiliconv

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

func TestParser_Fixture(t *testing.T) {
	html := mustReadTestHTML(t, "bilibili_video_sample.html")
	res, err := NewParser().Parse(context.Background(), &core.FetchResult{
		URL:  "https://www.bilibili.com/video/BV1X1Lj66ENe/",
		HTML: html,
	}, "output")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if strings.TrimSpace(res.Title) == "" {
		t.Fatalf("expected non-empty title")
	}
	if !strings.Contains(res.Title, "健身解剖") {
		t.Fatalf("unexpected title: %q", res.Title)
	}
	if len(res.Assets) == 0 {
		t.Fatalf("expected at least one asset")
	}
	hasVideo := false
	for _, a := range res.Assets {
		if a.Type == core.AssetTypeVideo {
			hasVideo = true
		}
	}
	if !hasVideo {
		t.Fatalf("expected a video asset")
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

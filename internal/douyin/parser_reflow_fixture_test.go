package douyinconv

import (
	"context"
	_ "embed"
	"strings"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

//go:embed testdata/douyin_reflow_sample.html
var douyinReflowSampleHTML string

// Real Douyin reflow/share HTML snapshot (video 7632652343653289266) for regression testing.
func TestParser_Parse_RealDouyinReflowFixture(t *testing.T) {
	if len(douyinReflowSampleHTML) < 1000 {
		t.Fatal("embedded fixture seems empty or truncated")
	}
	const pageURL = "https://www.iesdouyin.com/share/video/7632652343653289266"
	const wantTitle = "大家好！今天我们以夏天为主题，练习英语听力与理解能力。"
	const wantAccount = "英语家园Airt"

	res, err := NewParser().Parse(context.Background(), &core.FetchResult{
		URL:  pageURL,
		HTML: douyinReflowSampleHTML,
	}, "output")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := strings.TrimSpace(res.Title); got != wantTitle {
		t.Fatalf("title: want %q, got %q", wantTitle, got)
	}
	if got := strings.TrimSpace(res.AccountName); got != wantAccount {
		t.Fatalf("account: want %q, got %q", wantAccount, got)
	}
	if !strings.Contains(res.Markdown, wantTitle) {
		t.Fatalf("markdown should contain title, got len=%d", len(res.Markdown))
	}
	if !strings.Contains(res.Markdown, "作者: "+wantAccount) {
		t.Fatalf("markdown should list author, got:\n%s", res.Markdown)
	}

	var videoURLs []string
	for _, a := range res.Assets {
		if a.Type == core.AssetTypeVideo {
			videoURLs = append(videoURLs, a.SourceURL)
		}
	}
	if len(videoURLs) == 0 {
		t.Fatal("expected at least one video asset")
	}
	primary := videoURLs[0]
	if !strings.Contains(primary, "aweme.snssdk.com") || !strings.Contains(primary, "/play/") {
		t.Fatalf("expected no-watermark play URL, got %q", primary)
	}
	if strings.Contains(primary, "playwm") {
		t.Fatalf("expected playwm normalized to play, got %q", primary)
	}
}

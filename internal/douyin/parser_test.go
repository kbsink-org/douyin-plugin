package douyinconv

import (
	"context"
	"strings"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

func TestParser_Parse(t *testing.T) {
	html := `
<html>
  <head>
    <meta property="og:title" content="Douyin Test Video Fallback" />
    <meta name="author" content="Test Creator" />
    <meta property="og:image" content="https://cdn.example.com/cover.jpg" />
    <meta property="og:video" content="https://cdn.example.com/playwm.mp4" />
    <script>window.__INIT_PROPS__={"desc":"test:/\\*title?","play_addr":{"url_list":["https:\/\/cdn.example.com\/playwm\/abc.mp4"]}}</script>
  </head>
  <body>
    <video src="https://cdn.example.com/play.mp4"></video>
  </body>
</html>`
	res, err := NewParser().Parse(context.Background(), &core.FetchResult{
		URL:  "https://www.douyin.com/video/123",
		HTML: html,
	}, "output")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if strings.TrimSpace(res.Title) == "" {
		t.Fatalf("expected non-empty title")
	}
	if strings.TrimSpace(res.Markdown) == "" {
		t.Fatalf("expected non-empty markdown")
	}
	if res.Title != "test_____title_" {
		t.Fatalf("expected title from desc with sanitize, got %q", res.Title)
	}
	videoCount := 0
	hasNoWatermark := false
	for _, asset := range res.Assets {
		if asset.Type == core.AssetTypeVideo {
			videoCount++
			if strings.Contains(asset.SourceURL, "/play/") {
				hasNoWatermark = true
			}
		}
	}
	if videoCount == 0 {
		t.Fatalf("expected at least one video asset")
	}
	if !hasNoWatermark {
		t.Fatalf("expected at least one no-watermark play url")
	}
}

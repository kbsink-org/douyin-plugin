package douyin

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kbsink-org/kbsink/pkg/core"
)

var escapedURLPattern = regexp.MustCompile(`https?:\\?/\\?/[^"'\\\s>]+`)
var douyinPlayAddrPattern = regexp.MustCompile(`"play_addr"[^}]*"url_list"[^[]*\[\s*"([^"]+)"`)
var douyinDescPattern = regexp.MustCompile(`"desc"\s*:\s*"([^"]+)"`)
var douyinVideoIDInURLPattern = regexp.MustCompile(`(?i)(?:^|/)video/([^/?#]+)`)
var invalidTitleChars = regexp.MustCompile(`[\\/:*?"<>|]`)

// Parser extracts metadata/content/media from Douyin share HTML.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(_ context.Context, fetched *core.FetchResult, outputDir string) (*core.ArticleResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fetched.HTML))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	videoID := extractDouyinVideoIDFromURL(fetched.URL)
	playURL, titleFromHTML := extractVideoInfoFromHTML(fetched.HTML, videoID)
	title := firstNonEmpty(
		titleFromHTML,
		strings.TrimSpace(doc.Find(`meta[property="og:title"]`).AttrOr("content", "")),
		strings.TrimSpace(doc.Find(`meta[name="description"]`).AttrOr("content", "")),
		strings.TrimSpace(doc.Find("title").First().Text()),
	)
	account := firstNonEmpty(
		strings.TrimSpace(doc.Find(`meta[name="author"]`).AttrOr("content", "")),
		extractJSONField(fetched.HTML, "nickname"),
	)

	assets := make([]core.Asset, 0)
	seen := map[string]struct{}{}
	addAsset := func(assetType core.AssetType, src string) {
		src = strings.TrimSpace(src)
		if src == "" || !isFetchableAssetURL(src) {
			return
		}
		if _, ok := seen[src]; ok {
			return
		}
		seen[src] = struct{}{}
		assets = append(assets, core.Asset{
			Type:         assetType,
			SourceURL:    src,
			RelativePath: path.Join(outputDir, "images"),
		})
	}

	addAsset(core.AssetTypeVideo, playURL)
	addAsset(core.AssetTypeImage, doc.Find(`meta[property="og:image"]`).AttrOr("content", ""))
	addAsset(core.AssetTypeVideo, doc.Find(`meta[property="og:video"]`).AttrOr("content", ""))
	addAsset(core.AssetTypeVideo, doc.Find(`meta[property="og:video:url"]`).AttrOr("content", ""))
	addAsset(core.AssetTypeVideo, doc.Find(`meta[itemprop="contentUrl"]`).AttrOr("content", ""))
	addAsset(core.AssetTypeVideo, doc.Find(`video`).First().AttrOr("src", ""))
	doc.Find("video source").Each(func(_ int, sel *goquery.Selection) {
		addAsset(core.AssetTypeVideo, sel.AttrOr("src", ""))
	})

	for _, src := range extractURLsByExtensions(fetched.HTML, []string{".mp4", ".mov", ".m4v"}) {
		addAsset(core.AssetTypeVideo, src)
	}
	for _, src := range extractEscapedURLs(fetched.HTML) {
		if strings.Contains(strings.ToLower(src), ".mp4") {
			addAsset(core.AssetTypeVideo, src)
		}
	}
	md := strings.TrimSpace(title)
	if md != "" {
		md = "# " + md
	}
	if md == "" {
		md = "# Douyin"
	}
	if account != "" {
		md += "\n\n作者: " + account
	}
	md += buildVideoLinksMarkdown(assets)
	for _, asset := range assets {
		if asset.Type == core.AssetTypeImage {
			md += "\n\n![](" + asset.SourceURL + ")"
			break
		}
	}

	return &core.ArticleResult{
		Title:          title,
		AccountName:    account,
		SourceURL:      fetched.URL,
		Markdown:       md,
		Assets:         assets,
		Images:         imageAssetsFromGenericAssets(assets, outputDir),
		RawHTMLContent: fetched.HTML,
	}, nil
}

func extractEscapedURLs(raw string) []string {
	matches := escapedURLPattern.FindAllString(raw, -1)
	results := make([]string, 0, len(matches))
	for _, match := range matches {
		normalized := strings.ReplaceAll(match, `\/`, "/")
		if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") {
			results = append(results, normalized)
		}
	}
	return results
}

func extractVideoInfoFromHTML(html string, videoID string) (string, string) {
	videoURL := ""
	if m := douyinPlayAddrPattern.FindStringSubmatch(html); len(m) == 2 {
		videoURL = strings.TrimSpace(strings.ReplaceAll(m[1], `\/`, "/"))
		videoURL = strings.ReplaceAll(videoURL, "playwm", "play")
	}

	title := ""
	if m := douyinDescPattern.FindStringSubmatch(html); len(m) == 2 {
		title = strings.TrimSpace(strings.ReplaceAll(m[1], `\"`, `"`))
	}
	title = sanitizeDouyinTitle(title)
	if title == "" && strings.TrimSpace(videoID) != "" {
		title = "douyin_" + strings.TrimSpace(videoID)
	}
	if videoURL != "" {
		return videoURL, title
	}
	backupURL := "https://aweme.snssdk.com/aweme/v1/play/?video_id=" + url.QueryEscape(strings.TrimSpace(videoID))
	return backupURL, title
}

func sanitizeDouyinTitle(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = invalidTitleChars.ReplaceAllString(raw, "_")
	return strings.TrimSpace(raw)
}

func extractDouyinVideoIDFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	if m := douyinVideoIDInURLPattern.FindStringSubmatch(rawURL); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func extractURLsByExtensions(raw string, exts []string) []string {
	quotedURLPattern := regexp.MustCompile(`https?://[^"']+`)
	matches := quotedURLPattern.FindAllString(raw, -1)
	results := make([]string, 0)
	for _, match := range matches {
		lower := strings.ToLower(match)
		for _, ext := range exts {
			if strings.Contains(lower, ext) {
				results = append(results, strings.Split(match, `\u0026`)[0])
				break
			}
		}
	}
	return results
}

func buildVideoLinksMarkdown(assets []core.Asset) string {
	var b strings.Builder
	for _, asset := range assets {
		if asset.Type != core.AssetTypeVideo {
			continue
		}
		_, _ = b.WriteString("\n\n[video](" + asset.SourceURL + ")")
	}
	return b.String()
}

func extractJSONField(rawHTML, field string) string {
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(field) + `"\s*:\s*"([^"]+)"`)
	match := pattern.FindStringSubmatch(rawHTML)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func imageAssetsFromGenericAssets(assets []core.Asset, outputDir string) []core.ImageAsset {
	images := make([]core.ImageAsset, 0, len(assets))
	for _, asset := range assets {
		if asset.Type != core.AssetTypeImage {
			continue
		}
		images = append(images, core.ImageAsset{
			SourceURL:    asset.SourceURL,
			RelativePath: path.Join(outputDir, "images"),
		})
	}
	return images
}

func isFetchableAssetURL(src string) bool {
	lower := strings.ToLower(strings.TrimSpace(src))
	if lower == "" {
		return false
	}
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

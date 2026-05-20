package xhsconv

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kbsink-org/kbsink-plugins/internal/parseutil"
	"github.com/kbsink-org/kbsink/pkg/core"
)

// Parser extracts metadata/content/media from Xiaohongshu share HTML.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(_ context.Context, fetched *core.FetchResult, outputDir string) (*core.ArticleResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fetched.HTML))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	title := firstNonEmpty(
		strings.TrimSpace(doc.Find(`meta[property="og:title"]`).AttrOr("content", "")),
		strings.TrimSpace(doc.Find("title").First().Text()),
	)
	account := strings.TrimSpace(doc.Find(`meta[name="author"]`).AttrOr("content", ""))
	if account == "" {
		account = extractJSONField(fetched.HTML, "nickname")
	}

	contentSel := firstSelection(doc,
		"#detail-desc",
		".note-content",
		".content",
		"article",
	)

	assets := collectXhsMediaAssets(doc, contentSel, fetched.HTML, outputDir)
	md := ""
	rawHTML := fetched.HTML
	if contentSel != nil && contentSel.Length() > 0 {
		md = parseutil.SelectionToMarkdown(contentSel)
		if inner, htmlErr := contentSel.Html(); htmlErr == nil {
			rawHTML = inner
		}
	}

	md = strings.TrimSpace(md)
	if md == "" {
		md = fmt.Sprintf("# %s\n", title)
	}
	// Note text rarely includes carousel <img>; images live in Assets only. Emit markdown
	// lines with original CDN URLs so Converter can rewrite them to local paths.
	if imgBlock := buildImageLinksMarkdown(assets); imgBlock != "" {
		md = imgBlock + "\n\n" + md
	}
	md += buildVideoLinksMarkdown(assets)

	return &core.ArticleResult{
		Title:          title,
		AccountName:    account,
		SourceURL:      fetched.URL,
		Markdown:       md,
		Assets:         assets,
		Images:         imageAssetsFromGenericAssets(assets, outputDir),
		RawHTMLContent: rawHTML,
	}, nil
}

func firstSelection(doc *goquery.Document, selectors ...string) *goquery.Selection {
	for _, selector := range selectors {
		sel := doc.Find(selector).First()
		if sel.Length() > 0 {
			return sel
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func imgURLFromSelection(sel *goquery.Selection) string {
	for _, attr := range []string{
		"data-src", "src", "data-original", "data-url",
		"data-lazy-src", "data-preview-url", "data-orig", "data-original-src",
	} {
		if v := strings.TrimSpace(sel.AttrOr(attr, "")); v != "" {
			return v
		}
	}
	return ""
}

// collectXhsMediaAssets collects images/videos from note body, full-page img tags,
// extension-based URLs, and XHS CDN URLs embedded in JSON (including \u002F escapes).
func collectXhsMediaAssets(doc *goquery.Document, contentSel *goquery.Selection, rawHTML string, outputDir string) []core.Asset {
	assets := make([]core.Asset, 0)
	seen := map[string]struct{}{}
	addAsset := func(assetType core.AssetType, src string) {
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}
		if !parseutil.IsFetchableAssetURL(src) {
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
	addImage := func(src string) { addAsset(core.AssetTypeImage, src) }
	addImageIfNote := func(src string) {
		if isLikelyXSNoteImageURL(src) {
			addImage(src)
		}
	}

	if contentSel != nil {
		contentSel.Find("img").Each(func(_ int, sel *goquery.Selection) {
			addImage(firstNonEmpty(
				imgURLFromSelection(sel),
			))
		})
		contentSel.Find("video source, video").Each(func(_ int, sel *goquery.Selection) {
			src := strings.TrimSpace(sel.AttrOr("src", ""))
			addAsset(core.AssetTypeVideo, src)
		})
	}

	for _, src := range extractURLsByExtensions(rawHTML, []string{".mp4", ".mov", ".m4v"}) {
		addAsset(core.AssetTypeVideo, src)
	}
	for _, src := range extractURLsByExtensions(rawHTML, []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}) {
		addImage(src)
	}

	norm := normalizeHTMLJSONStringEscapes(rawHTML)
	for _, src := range extractURLsByExtensions(norm, []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}) {
		addImage(src)
	}
	for _, src := range extractXHSCDNImageURLs(norm) {
		addImageIfNote(src)
	}

	// Whole-document imgs (carousel / media column lives outside .note-content).
	doc.Find("img").Each(func(_ int, sel *goquery.Selection) {
		u := imgURLFromSelection(sel)
		if u == "" {
			return
		}
		addImageIfNote(u)
	})

	return assets
}

// normalizeHTMLJSONStringEscapes turns common JSON-escaped sequences in inline HTML/JSON
// (e.g. http:\u002F\u002Fsns-webpic..., https:\/\/host\/path) into plain URLs so regexes match full paths.
func normalizeHTMLJSONStringEscapes(s string) string {
	// JSON string escapes for "/" — without this, CDN regexes stop at "\" and yield a host-only URL.
	s = strings.ReplaceAll(s, `\/`, `/`)
	for _, r := range []struct{ from, to string }{
		{`\u002f`, `/`},
		{`\u002F`, `/`},
		{`\u003a`, `:`},
		{`\u003A`, `:`},
	} {
		s = strings.ReplaceAll(s, r.from, r.to)
	}
	return s
}

var (
	reSNSWebpicURL = regexp.MustCompile(`(?i)https?://sns-webpic[^"'\s<>)\\]+`)
	reXHSNotePic   = regexp.MustCompile(`(?i)https?://[a-z0-9.-]*xhscdn\.com/[^"'\s<>)\\]+`)
)

func extractXHSCDNImageURLs(normHTML string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)
	collect := func(matches []string) {
		for _, m := range matches {
			u := trimURLSuffixNoise(m)
			if u == "" {
				continue
			}
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			out = append(out, u)
		}
	}
	collect(reSNSWebpicURL.FindAllString(normHTML, -1))
	for _, m := range reXHSNotePic.FindAllString(normHTML, -1) {
		u := trimURLSuffixNoise(m)
		if isLikelyXSNoteImageURL(u) {
			collect([]string{u})
		}
	}
	return out
}

func trimURLSuffixNoise(u string) string {
	u = strings.TrimSpace(u)
	for _, suf := range []string{`\`, `",`, `"`, `',`, `'`, `)`, `]`, `}`, `,`} {
		u = strings.TrimSuffix(u, suf)
	}
	return strings.TrimSpace(u)
}

// isLikelyXSNoteImageURL filters CDN hits to note/gallery images, skipping avatars and tiny UI assets.
func isLikelyXSNoteImageURL(src string) bool {
	u := strings.TrimSpace(strings.ToLower(src))
	if u == "" {
		return false
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return false
	}
	if strings.Contains(u, "/avatar/") {
		return false
	}
	if strings.Contains(u, "emoji") || strings.Contains(u, "icon") {
		return false
	}
	if strings.Contains(u, "sns-webpic") {
		parsed, err := url.Parse(src)
		if err != nil || strings.Trim(parsed.Path, "/") == "" {
			return false
		}
		return true
	}
	if strings.Contains(u, "xhscdn.com") {
		parsed, err := url.Parse(src)
		if err != nil || strings.Trim(parsed.Path, "/") == "" {
			return false
		}
		if strings.Contains(u, "1040g") {
			return true
		}
		if strings.Contains(u, "!nd_") || strings.Contains(u, "!wb_") || strings.Contains(u, "!prv_") {
			return true
		}
		if strings.Contains(u, "/pic/") || strings.Contains(u, "/notes/") {
			return true
		}
	}
	return false
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

func buildImageLinksMarkdown(assets []core.Asset) string {
	var b strings.Builder
	for _, asset := range assets {
		if asset.Type != core.AssetTypeImage {
			continue
		}
		if b.Len() > 0 {
			_, _ = b.WriteString("\n")
		}
		_, _ = b.WriteString("![](")
		_, _ = b.WriteString(asset.SourceURL)
		_, _ = b.WriteString(")")
	}
	return b.String()
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


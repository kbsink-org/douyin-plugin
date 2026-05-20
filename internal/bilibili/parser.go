package bilibiliconv

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kbsink-org/kbsink-plugins/internal/parseutil"
	"github.com/kbsink-org/kbsink/pkg/core"
)

var (
	biliJSONTitlePattern  = regexp.MustCompile(`"title"\s*:\s*"((?:\\.|[^"\\])*)"`)
	biliJSONPicPattern    = regexp.MustCompile(`"pic"\s*:\s*"(https?://[^"\\]+)"`)
	biliJSONBaseURLPattern = regexp.MustCompile(`"baseUrl"\s*:\s*"(https?://[^"\\]+)"`)
	biliJSONDurlPattern   = regexp.MustCompile(`"durl"\s*:\s*\[\s*\{[^}]*"url"\s*:\s*"(https?://[^"\\]+)"`)
	biliTitleSuffixPattern = regexp.MustCompile(`_\s*哔哩哔哩[_\s-]*bilibili\s*$`)
)

// Parser extracts metadata/content/media from Bilibili video HTML.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(_ context.Context, fetched *core.FetchResult, outputDir string) (*core.ArticleResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fetched.HTML))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	title := parseutil.FirstNonEmpty(
		cleanBilibiliTitle(doc.Find(`meta[property="og:title"]`).AttrOr("content", "")),
		cleanBilibiliTitle(extractJSONUnescaped(fetched.HTML, biliJSONTitlePattern)),
		cleanBilibiliTitle(doc.Find("title").First().Text()),
	)
	account := parseutil.FirstNonEmpty(
		strings.TrimSpace(doc.Find(`meta[name="author"]`).AttrOr("content", "")),
		extractOwnerName(fetched.HTML),
	)
	desc := parseutil.FirstNonEmpty(
		strings.TrimSpace(doc.Find(`meta[property="og:description"]`).AttrOr("content", "")),
		strings.TrimSpace(doc.Find(`meta[name="description"]`).AttrOr("content", "")),
	)

	assets := collectBilibiliAssets(doc, fetched.HTML, outputDir)

	md := strings.TrimSpace(title)
	if md != "" {
		md = "# " + md
	} else {
		md = "# Bilibili"
	}
	if account != "" {
		md += "\n\nUP主: " + account
	}
	if desc != "" {
		md += "\n\n" + desc
	}
	md += parseutil.BuildVideoLinksMarkdown(assets)
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
		Images:         parseutil.ImageAssetsFromGenericAssets(assets, outputDir),
		RawHTMLContent: fetched.HTML,
	}, nil
}

func collectBilibiliAssets(doc *goquery.Document, rawHTML, outputDir string) []core.Asset {
	assets := make([]core.Asset, 0)
	seen := map[string]struct{}{}
	add := func(assetType core.AssetType, src string) {
		src = strings.TrimSpace(src)
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

	add(core.AssetTypeImage, doc.Find(`meta[property="og:image"]`).AttrOr("content", ""))
	if m := biliJSONPicPattern.FindStringSubmatch(rawHTML); len(m) == 2 {
		add(core.AssetTypeImage, unescapeJSONString(m[1]))
	}

	add(core.AssetTypeVideo, doc.Find(`meta[property="og:video"]`).AttrOr("content", ""))
	add(core.AssetTypeVideo, doc.Find(`meta[property="og:video:url"]`).AttrOr("content", ""))
	if m := biliJSONBaseURLPattern.FindStringSubmatch(rawHTML); len(m) == 2 {
		add(core.AssetTypeVideo, unescapeJSONString(m[1]))
	}
	if m := biliJSONDurlPattern.FindStringSubmatch(rawHTML); len(m) == 2 {
		add(core.AssetTypeVideo, unescapeJSONString(m[1]))
	}
	for _, src := range extractURLsByExtensions(rawHTML, []string{".mp4", ".m4v", ".flv"}) {
		add(core.AssetTypeVideo, src)
	}
	return assets
}

func cleanBilibiliTitle(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = biliTitleSuffixPattern.ReplaceAllString(raw, "")
	return strings.TrimSpace(raw)
}

func extractOwnerName(rawHTML string) string {
	pattern := regexp.MustCompile(`"owner"\s*:\s*\{[^}]*"name"\s*:\s*"([^"]+)"`)
	if m := pattern.FindStringSubmatch(rawHTML); len(m) == 2 {
		return strings.TrimSpace(unescapeJSONString(m[1]))
	}
	return ""
}

func extractJSONUnescaped(raw string, re *regexp.Regexp) string {
	if m := re.FindStringSubmatch(raw); len(m) == 2 {
		return unescapeJSONString(m[1])
	}
	return ""
}

func unescapeJSONString(s string) string {
	s = strings.ReplaceAll(s, `\/`, "/")
	s = strings.ReplaceAll(s, `\u002F`, "/")
	s = strings.ReplaceAll(s, `\u002f`, "/")
	s = strings.ReplaceAll(s, `\u0026`, "&")
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	return strings.TrimSpace(s)
}

func extractURLsByExtensions(raw string, exts []string) []string {
	quotedURLPattern := regexp.MustCompile(`https?://[^"'\s<>\\]+`)
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

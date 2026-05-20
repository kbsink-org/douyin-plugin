package zhihuconv

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

var zhihuTitleSuffixPattern = regexp.MustCompile(`\s*-\s*知乎\s*$`)

// Parser extracts metadata/content from Zhihu zhuanlan article HTML.
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
		cleanZhihuTitle(doc.Find(`meta[property="og:title"]`).AttrOr("content", "")),
		cleanZhihuTitle(doc.Find("h1.Post-Title").First().Text()),
		cleanZhihuTitle(doc.Find("h1").First().Text()),
		cleanZhihuTitle(doc.Find("title").First().Text()),
	)
	account := parseutil.FirstNonEmpty(
		strings.TrimSpace(doc.Find(`meta[name="author"]`).AttrOr("content", "")),
		strings.TrimSpace(doc.Find(".AuthorInfo-name").First().Text()),
		strings.TrimSpace(doc.Find(`a.UserLink-link`).First().Text()),
	)

	contentSel := firstSelection(doc,
		".Post-RichText",
		".RichText.ztext",
		"article .RichText",
		"article",
		".ArticleItem-content",
	)

	assets := collectZhihuAssets(doc, contentSel, outputDir)
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
		if desc := strings.TrimSpace(doc.Find(`meta[property="og:description"]`).AttrOr("content", "")); desc != "" {
			md = desc
		}
	}
	if md == "" && title != "" {
		md = "# " + title
	} else if title != "" && !strings.HasPrefix(md, "# ") {
		md = "# " + title + "\n\n" + md
	}
	if imgBlock := buildImageLinksMarkdown(assets); imgBlock != "" && !strings.Contains(md, imgBlock) {
		md = imgBlock + "\n\n" + md
	}

	return &core.ArticleResult{
		Title:          title,
		AccountName:    account,
		SourceURL:      fetched.URL,
		Markdown:       md,
		Assets:         assets,
		Images:         parseutil.ImageAssetsFromGenericAssets(assets, outputDir),
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

func collectZhihuAssets(doc *goquery.Document, contentSel *goquery.Selection, outputDir string) []core.Asset {
	assets := make([]core.Asset, 0)
	seen := map[string]struct{}{}
	addImage := func(src string) {
		src = strings.TrimSpace(src)
		if !parseutil.IsFetchableAssetURL(src) {
			return
		}
		if _, ok := seen[src]; ok {
			return
		}
		seen[src] = struct{}{}
		assets = append(assets, core.Asset{
			Type:         core.AssetTypeImage,
			SourceURL:    src,
			RelativePath: path.Join(outputDir, "images"),
		})
	}

	addImage(doc.Find(`meta[property="og:image"]`).AttrOr("content", ""))
	if contentSel != nil {
		contentSel.Find("img").Each(func(_ int, sel *goquery.Selection) {
			src := parseutil.FirstNonEmpty(
				sel.AttrOr("data-original", ""),
				sel.AttrOr("data-actualsrc", ""),
				sel.AttrOr("data-src", ""),
				sel.AttrOr("src", ""),
			)
			addImage(src)
		})
	}
	return assets
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

func cleanZhihuTitle(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = zhihuTitleSuffixPattern.ReplaceAllString(raw, "")
	return strings.TrimSpace(raw)
}

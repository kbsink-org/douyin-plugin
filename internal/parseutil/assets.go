package parseutil

import (
	"path"
	"strings"

	"github.com/kbsink-org/kbsink/pkg/core"
)

// BuildVideoLinksMarkdown appends [video](url) lines for video assets.
func BuildVideoLinksMarkdown(assets []core.Asset) string {
	var b strings.Builder
	for _, asset := range assets {
		if asset.Type != core.AssetTypeVideo {
			continue
		}
		_, _ = b.WriteString("\n\n[video](" + asset.SourceURL + ")")
	}
	return b.String()
}

// ImageAssetsFromGenericAssets maps image-type assets to ImageAsset (compatibility field).
func ImageAssetsFromGenericAssets(assets []core.Asset, outputDir string) []core.ImageAsset {
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

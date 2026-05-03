package main

import (
	"encoding/json"
	"io"

	"github.com/kbsink-org/kbsink/pkg/core"
)

type cliJSONResult struct {
	OK           bool   `json:"ok"`
	Error        string `json:"error,omitempty"`
	SourceURL    string `json:"source_url,omitempty"`
	Title        string `json:"title,omitempty"`
	MarkdownPath string `json:"markdown_path,omitempty"`
	OutputDir    string `json:"output_dir,omitempty"`
	Markdown     string `json:"markdown,omitempty"`
	ImageCount   int    `json:"images,omitempty"`
	VideoCount   int    `json:"videos,omitempty"`
}

func writeCLIJSON(w io.Writer, srcURL string, res *core.ArticleResult, err error, includeMarkdown bool) error {
	out := cliJSONResult{SourceURL: srcURL}
	if err != nil {
		out.OK = false
		out.Error = err.Error()
	} else if res != nil {
		out.OK = true
		out.Title = res.Title
		out.MarkdownPath = res.MarkdownPath
		out.OutputDir = res.OutputDir
		if includeMarkdown {
			out.Markdown = res.Markdown
		}
		out.ImageCount = len(res.Images)
		for _, a := range res.Assets {
			if a.Type == core.AssetTypeVideo {
				out.VideoCount++
			}
		}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(out)
}

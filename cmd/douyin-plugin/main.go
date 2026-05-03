package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	kbsink "github.com/kbsink-org/kbsink/pkg"
	"github.com/kbsink-org/kbsink/pkg/core"
	"github.com/kbsink-org/douyin-plugin/pkg/douyin"
)

func main() {
	var (
		outputRoot = flag.String("o", "output", "output root directory")
		timeout    = flag.Duration("timeout", 60*time.Second, "timeout for fetch+parse")
		printOnly  = flag.Bool("print", false, "print markdown to stdout (still downloads assets when not print-only)")
		videoMode  = flag.String("video-mode", "link", "video markdown mode: link|embed")
		format     = flag.String("format", "text", "output format: text|json (for full CLI with --plugin douyin use kb-sink-md-douyin)")
	)
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  %s [flags] <douyin-share-url-or-text>\n\n"+
			"Runs kbsink Converter with Douyin driver+parser.\n"+
			"For the same flags as kb-sink-md with Douyin, use the kb-sink-md-douyin binary.\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	raw := flag.Arg(0)

	ctx := context.Background()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	mode, err := resolveVideoMode(*videoMode)
	if err != nil {
		emitErr(*format, raw, err)
		os.Exit(1)
	}

	converter := kbsink.NewConverter(
		kbsink.WithParser(douyin.NewParser()),
		kbsink.WithDriver(douyin.NewDriver(nil)),
	)

	res, err := converter.Convert(ctx, raw, core.ConvertOptions{
		OutputRoot: *outputRoot,
		VideoMode:  mode,
	})
	if err != nil {
		emitErr(*format, raw, err)
		os.Exit(1)
	}

	if err := emitOK(*format, raw, res, *printOnly); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "emit output: %v\n", err)
		os.Exit(1)
	}
}

func emitErr(format, src string, err error) {
	if strings.TrimSpace(strings.ToLower(format)) == "json" {
		_ = writeCLIJSON(os.Stdout, src, nil, err, false)
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
}

func emitOK(format, src string, res *core.ArticleResult, printOnly bool) error {
	if strings.TrimSpace(strings.ToLower(format)) == "json" {
		return writeCLIJSON(os.Stdout, src, res, nil, printOnly)
	}
	if printOnly {
		_, _ = fmt.Fprint(os.Stdout, res.Markdown)
		return nil
	}
	_, _ = fmt.Fprintf(os.Stdout, "title: %s\n", res.Title)
	_, _ = fmt.Fprintf(os.Stdout, "markdown: %s\n", res.MarkdownPath)
	_, _ = fmt.Fprintf(os.Stdout, "images: %d\n", len(res.Images))
	videoCount := 0
	for _, asset := range res.Assets {
		if asset.Type == core.AssetTypeVideo {
			videoCount++
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, "videos: %d\n", videoCount)
	return nil
}

func resolveVideoMode(raw string) (core.VideoMode, error) {
	mode := strings.TrimSpace(strings.ToLower(raw))
	switch mode {
	case "", string(core.VideoModeLink):
		return core.VideoModeLink, nil
	case string(core.VideoModeEmbed):
		return core.VideoModeEmbed, nil
	default:
		return "", fmt.Errorf("unsupported video mode %q, expected link|embed", raw)
	}
}

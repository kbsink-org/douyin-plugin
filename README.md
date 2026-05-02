# douyin-plugin

English | [简体中文](README.zh-CN.md)

A [kbsink](https://github.com/kbsink-org/kbsink) plugin for Douyin share links: a `Driver` (fetch page HTML) and a `Parser` (turn HTML into Markdown) for use with `Converter`. Includes `cmd/douyin-plugin` for local integration and smoke testing.

## Requirements

- Go 1.26 or newer (see `go.mod`)

## Library usage

```go
import (
    kbsink "github.com/kbsink-org/kbsink/pkg"
    "github.com/kbsink-org/kbsink/pkg/core"
    "github.com/kbsink-org/douyin-plugin/pkg/douyin"
)

converter := kbsink.NewConverter(
    kbsink.WithParser(douyin.NewParser()),
    kbsink.WithDriver(douyin.NewDriver(nil)), // or pass a custom *http.Client
)

res, err := converter.Convert(ctx, shareURLOrShareText, core.ConvertOptions{
    OutputRoot: "output",
    VideoMode:  core.VideoModeLink, // or core.VideoModeEmbed
})
```

Input may be a bare share URL or share text containing a URL (the first `http(s)` URL in the text is used).

## Command-line

From the repository root:

```bash
go run ./cmd/douyin-plugin --help
```

Examples:

```bash
# Write assets under ./output and print a short summary to the terminal
go run ./cmd/douyin-plugin -o output "https://v.douyin.com/xxxx/"

# Print Markdown to stdout only (asset handling still follows Converter unless you change options)
go run ./cmd/douyin-plugin -print "https://v.douyin.com/xxxx/"

# Video in Markdown as link or embed
go run ./cmd/douyin-plugin -video-mode=link "..."
go run ./cmd/douyin-plugin -video-mode=embed "..."

# Fetch+parse timeout (default 60s; 0 means no limit)
go run ./cmd/douyin-plugin -timeout=90s "..."
```

## Tests

```bash
go test ./...
```

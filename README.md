# douyin-plugin

English | [简体中文](README.zh-CN.md)

A [kbsink](https://github.com/kbsink-org/kbsink) plugin for Douyin share links: a `Driver` (fetch page HTML) and a `Parser` (turn HTML into Markdown) for use with `Converter`. Includes `cmd/douyin-plugin` for local integration and smoke testing.

## Requirements

- Go **1.24** or newer, aligned with [kbsink](https://github.com/kbsink-org/kbsink) (see `go.mod`).
- Local development uses `replace github.com/kbsink-org/kbsink => ../kbsink` when this repo sits next to `kbsink` in the same parent folder.

## Releases / binaries

Pushing a tag `v*` runs GitHub Actions: cross-builds **`kb-sink-md-douyin`** (same CLI as `kb-sink-md`, with Douyin Parser+Driver registered for `--plugin douyin`; `kbsink.Converter` runs in-process) plus optional **`douyin-plugin`** smoke binary. Archives and `SHA256SUMS.txt` attach to the GitHub Release. CI checks out `kbsink-org/kbsink` next to this module so the `replace` directive resolves on the runner.

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

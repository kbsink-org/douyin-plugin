# douyin-plugin

English | [简体中文](README.zh-CN.md)

Go module that adds **Douyin** (share links / share text) support to [kbsink](https://github.com/kbsink-org/kbsink): a **`Parser`** (HTML → article structure / markdown) and a **`Driver`** (HTTP fetch). Import **`github.com/kbsink-org/douyin-plugin/pkg/douyin`**.

**Command line:** use [kbsink-cli](https://github.com/kbsink-org/kbsink-cli) for a ready-made `kbsink` binary (WeChat, Xiaohongshu, Douyin).

## Requirements

- **Go 1.25+** (see `go.mod`; stay aligned with the kbsink version you depend on).
- Local dev: a [workspace](https://go.dev/ref/mod#workspaces) that includes this repo and `kbsink`, or in `go.mod`:

  ```text
  replace github.com/kbsink-org/kbsink => ../kbsink
  ```

## Use with `kbsink.Converter`

```go
import (
	kbsink "github.com/kbsink-org/kbsink/pkg"
	"github.com/kbsink-org/kbsink/pkg/core"
	"github.com/kbsink-org/douyin-plugin/pkg/douyin"
)

converter := kbsink.NewConverter(
	kbsink.WithParser(douyin.NewParser()),
	kbsink.WithDriver(douyin.NewDriver(nil)), // or your own *http.Client
)

res, err := converter.Convert(ctx, shareURLOrShareText, core.ConvertOptions{
	OutputRoot: "output",
	VideoMode:  core.VideoModeLink, // or core.VideoModeEmbed
})
```

`Convert` accepts a bare `https://v.douyin.com/...` URL or arbitrary text; the **first** `http(s)` URL in the string is used.

## Optional: `pluginreg` by name

`pkg/douyin` only exports **`NewParser`** and **`NewDriver`**. It does **not** ship a `core.Plugin` type. If you want `pluginreg.Lookup("douyin")` style registration, implement `core.Plugin` yourself. The maintained reference adapter is **[kbsink-cli/internal/plugin/douyin](https://github.com/kbsink-org/kbsink-cli/tree/main/internal/plugin/douyin)** (`douyin.New()` + `pluginreg.Register`); the shape is:

```go
type douyinPlugin struct{}

func (douyinPlugin) Name() string { return "douyin" }

func (douyinPlugin) NewComponents(c *http.Client) (core.Parser, core.Driver, error) {
	return douyin.NewParser(), douyin.NewDriver(c), nil
}
```

## Releases and binaries

- Tags **`v*`** mark **Go module** releases for `go get`.
- This repo does **not** publish its own CLI binaries; prebuilt **`kbsink`** builds and multi-platform archives live under **[kbsink-cli releases](https://github.com/kbsink-org/kbsink-cli/releases)**.

## Tests

```bash
go test ./...
```

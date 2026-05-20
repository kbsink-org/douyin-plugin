# kbsink-plugins

English | [简体中文](README.zh-CN.md)

Go module of **platform plugins** for [kbsink](https://github.com/kbsink-org/kbsink): each package under `pkg/<name>` exports a **`Parser`** (HTML → article structure / markdown) and a **`Driver`** (HTTP fetch).

| Package | Platform |
|---------|----------|
| `pkg/wechat` | WeChat Official Account articles |
| `pkg/xhs` | Xiaohongshu (小红书) share links |
| `pkg/douyin` | Douyin share links / share text |
| `pkg/bilibili` | Bilibili video pages (`BV…` URLs, `b23.tv`) |
| `pkg/zhihu` | Zhihu zhuanlan articles (`zhuanlan.zhihu.com/p/…`) |

**CLI:** use [kbsink-cli](https://github.com/kbsink-org/kbsink-cli) for a ready-made `kbsink` binary with all plugins wired.

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
	"github.com/kbsink-org/kbsink-plugins/pkg/wechat"
)

converter := kbsink.NewConverter(
	kbsink.WithParser(wechat.NewParser()),
	kbsink.WithDriver(wechat.NewDriver(nil)),
)

res, err := converter.Convert(ctx, articleURL, core.ConvertOptions{
	OutputRoot: "output",
})
```

## Optional: `pluginreg` by name

Each `pkg/<name>` exports **`NewParser`** and **`NewDriver`** only. For `pluginreg.Lookup("wechat")` style registration, see **[kbsink-cli/internal/plugin](https://github.com/kbsink-org/kbsink-cli/tree/main/internal/plugin)**.

## Tests

```bash
go test ./...
```

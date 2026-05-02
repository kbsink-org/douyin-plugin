# douyin-plugin

[English](README.md) | 简体中文

[kbsink](https://github.com/kbsink-org/kbsink) 的抖音分享链接插件：提供 `Driver`（拉取页面 HTML）与 `Parser`（解析为 Markdown），供 `Converter` 组合使用。附带 `cmd/douyin-plugin` 命令行，用于本地联调与冒烟测试。

## 要求

- Go 1.26+（见 `go.mod`）

## 作为库使用

```go
import (
    kbsink "github.com/kbsink-org/kbsink/pkg"
    "github.com/kbsink-org/kbsink/pkg/core"
    "github.com/kbsink-org/douyin-plugin/pkg/douyin"
)

converter := kbsink.NewConverter(
    kbsink.WithParser(douyin.NewParser()),
    kbsink.WithDriver(douyin.NewDriver(nil)), // 或传入自定义 *http.Client
)

res, err := converter.Convert(ctx, shareURLOrShareText, core.ConvertOptions{
    OutputRoot: "output",
    VideoMode:  core.VideoModeLink, // 或 core.VideoModeEmbed
})
```

输入可以是纯分享 URL，或包含 URL 的分享文案（会从文本中提取第一个 `http(s)` 链接）。

## 命令行

在项目根目录：

```bash
go run ./cmd/douyin-plugin --help
```

示例：

```bash
# 写出资源到 ./output，并在终端打印摘要
go run ./cmd/douyin-plugin -o output "https://v.douyin.com/xxxx/"

# 只把 Markdown 打到 stdout（仍会按 Converter 行为处理资源，除非你再配合其它选项）
go run ./cmd/douyin-plugin -print "https://v.douyin.com/xxxx/"

# 视频在 Markdown 中写作链接或 embed
go run ./cmd/douyin-plugin -video-mode=link "..."
go run ./cmd/douyin-plugin -video-mode=embed "..."

# 拉取+解析超时（默认 60s，0 表示不限制）
go run ./cmd/douyin-plugin -timeout=90s "..."
```

## 测试

```bash
go test ./...
```

# kbsink-plugins

[English](README.md) | 简体中文

面向 [kbsink](https://github.com/kbsink-org/kbsink) 的**平台插件** Go 模块：`pkg/<name>` 下每个包提供 **`Parser`**（HTML → 文章结构 / Markdown）与 **`Driver`**（HTTP 拉取）。

| 包 | 平台 |
|----|------|
| `pkg/wechat` | 微信公众号文章 |
| `pkg/xhs` | 小红书分享链接 |
| `pkg/douyin` | 抖音分享链接 / 分享文案 |
| `pkg/bilibili` | B 站视频页（`BV…` 链接、`b23.tv`） |
| `pkg/zhihu` | 知乎专栏文章（`zhuanlan.zhihu.com/p/…`） |

**命令行：** 使用 [kbsink-cli](https://github.com/kbsink-org/kbsink-cli) 获取已集成全部插件的 `kbsink` 二进制。

## 依赖

- **Go 1.25+**（见 `go.mod`；与所依赖的 kbsink 版本保持一致）。
- 本地开发：用 [workspace](https://go.dev/ref/mod#workspaces) 同时包含本仓库与 `kbsink`，或在 `go.mod` 中：

  ```text
  replace github.com/kbsink-org/kbsink => ../kbsink
  ```

## 与 `kbsink.Converter` 配合

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

## 测试

```bash
go test ./...
```

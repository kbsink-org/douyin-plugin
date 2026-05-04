# douyin-plugin

[English](README.md) | 简体中文

本仓库是面向 [kbsink](https://github.com/kbsink-org/kbsink) 的 **Go 库**：为抖音分享链接 / 分享文案提供 **`Parser`**（HTML → 文章结构 / Markdown）与 **`Driver`**（HTTP 拉取页面）。在业务代码中导入 **`github.com/kbsink-org/douyin-plugin/pkg/douyin`**。

**命令行：** 直接使用 [kbsink-cli](https://github.com/kbsink-org/kbsink-cli) 打包好的 `kbsink`（微信、小红书、抖音）。

## 要求

- **Go 1.25+**（以 `go.mod` 为准；与所依赖的 kbsink 版本保持兼容）。
- 本地开发：用 [Go workspace](https://go.dev/ref/mod#workspaces) 同时包含本仓库与 `kbsink`，或在 `go.mod` 中写：

  ```text
  replace github.com/kbsink-org/kbsink => ../kbsink
  ```

## 与 `kbsink.Converter` 集成

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

`Convert` 既可传裸的 `https://v.douyin.com/...` 链接，也可传整段文案；会取字符串中 **第一个** `http(s)` URL 作为输入。

## 可选：按名称走 `pluginreg`

`pkg/douyin` 只导出 **`NewParser`**、**`NewDriver`**，**不包含** `core.Plugin` 类型。若需要 `pluginreg.Lookup("douyin")` 这类按名解析，请自行实现 `core.Plugin`。推荐对照 **[kbsink-cli/internal/plugin/douyin](https://github.com/kbsink-org/kbsink-cli/tree/main/internal/plugin/douyin)**（`douyin.New()` + `pluginreg.Register`），核心形态如下：

```go
type douyinPlugin struct{}

func (douyinPlugin) Name() string { return "douyin" }

func (douyinPlugin) NewComponents(c *http.Client) (core.Parser, core.Driver, error) {
	return douyin.NewParser(), douyin.NewDriver(c), nil
}
```

## 发版与二进制

- **`v*`** 标签用于 **Go 模块** 版本（`go get`）。
- 本仓库 **不发布** 独立 CLI；预编译 **`kbsink`** 与多平台包见 **[kbsink-cli Releases](https://github.com/kbsink-org/kbsink-cli/releases)**。

## 测试

```bash
go test ./...
```

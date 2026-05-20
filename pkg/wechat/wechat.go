// Package wechat provides the WeChat article parser and driver for kbsink.
package wechat

import (
	wconv "github.com/kbsink-org/kbsink-plugins/internal/wechat"
	"github.com/kbsink-org/kbsink/pkg/core"
	"github.com/kbsink-org/kbsink/pkg/driver"
)

const userAgent = "Mozilla/5.0 (compatible; wechatmd/1.0)"

// NewParser returns a WeChat HTML parser for kbsink.Converter.
func NewParser() core.Parser { return wconv.NewParser() }

// NewDriver returns an HTTP fetch driver for WeChat article URLs.
func NewDriver(c core.HTTPClient) core.Driver { return driver.NewHTMLDriver(c, userAgent, nil) }

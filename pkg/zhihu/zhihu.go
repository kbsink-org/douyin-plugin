// Package zhihu provides the Zhihu zhuanlan article parser and driver for kbsink.
package zhihu

import (
	zconv "github.com/kbsink-org/kbsink-plugins/internal/zhihu"
	"github.com/kbsink-org/kbsink/pkg/core"
)

// NewParser returns a Zhihu article parser for kbsink.Converter.
func NewParser() core.Parser { return zconv.NewParser() }

// NewDriver returns an HTTP fetch driver for Zhihu article URLs.
func NewDriver(c core.HTTPClient) core.Driver { return zconv.NewHTTPDriver(c) }

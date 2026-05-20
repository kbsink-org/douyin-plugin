// Package xhs provides the Xiaohongshu (小红书) parser and driver for kbsink.
package xhs

import (
	xconv "github.com/kbsink-org/kbsink-plugins/internal/xhs"
	"github.com/kbsink-org/kbsink/pkg/core"
	"github.com/kbsink-org/kbsink/pkg/driver"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"

// NewParser returns an XHS HTML parser for kbsink.Converter.
func NewParser() core.Parser { return xconv.NewParser() }

// NewDriver returns an HTTP fetch driver for Xiaohongshu share/note URLs.
func NewDriver(c core.HTTPClient) core.Driver { return driver.NewHTMLDriver(c, userAgent, nil) }

// DefaultUserAgent is the User-Agent used by NewDriver.
func DefaultUserAgent() string { return userAgent }

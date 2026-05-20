// Package douyin provides the Douyin parser and driver for kbsink.
// Implementation lives in internal/douyin (package douyinconv); import this package from other modules (e.g. kbsink-server).
package douyin

import (
	"net/http"

	dconv "github.com/kbsink-org/kbsink-plugins/internal/douyin"
	"github.com/kbsink-org/kbsink/pkg/core"
)

// NewParser returns a Douyin HTML parser for kbsink.Converter.
func NewParser() core.Parser { return dconv.NewParser() }

// NewDriver returns a fetch driver for Douyin share URLs or share text.
func NewDriver(c *http.Client) core.Driver { return dconv.NewHTTPDriver(c) }

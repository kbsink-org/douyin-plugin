// Package bilibili provides the Bilibili video parser and driver for kbsink.
package bilibili

import (
	"net/http"

	bconv "github.com/kbsink-org/kbsink-plugins/internal/bilibili"
	"github.com/kbsink-org/kbsink/pkg/core"
)

// NewParser returns a Bilibili HTML parser for kbsink.Converter.
func NewParser() core.Parser { return bconv.NewParser() }

// NewDriver returns an HTTP fetch driver for Bilibili video URLs or share text.
func NewDriver(c *http.Client) core.Driver { return bconv.NewHTTPDriver(c) }

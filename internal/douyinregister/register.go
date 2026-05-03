// Package douyinregister registers Douyin with kbsink's in-process plugin registry.
// Blank-import only from CLI entrypoints (e.g. kb-sink-md-douyin); library code should use pkg/douyin directly.
package douyinregister

import (
	"github.com/kbsink-org/douyin-plugin/pkg/douyin"
	"github.com/kbsink-org/kbsink/pkg/pluginreg"
)

func init() {
	pluginreg.Register(douyin.Plugin{})
}

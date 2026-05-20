package parseutil

import "strings"

// IsFetchableAssetURL reports whether src is an absolute http(s) URL suitable for download.
func IsFetchableAssetURL(src string) bool {
	lower := strings.ToLower(strings.TrimSpace(src))
	if lower == "" {
		return false
	}
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

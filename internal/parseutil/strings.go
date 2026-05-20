package parseutil

import "strings"

// FirstNonEmpty returns the first non-blank string.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// ExtractShareURL returns the first http(s) URL in raw text (trimmed of trailing punctuation).
func ExtractShareURL(rawInput string) string {
	rawInput = strings.TrimSpace(rawInput)
	if rawInput == "" {
		return ""
	}
	const prefix = "http"
	idx := strings.Index(strings.ToLower(rawInput), prefix)
	if idx < 0 {
		return ""
	}
	rest := rawInput[idx:]
	end := len(rest)
	for i, r := range rest {
		switch r {
		case ' ', '\n', '\t', '"', '\'', '>', '）', '】', '》', '」':
			end = i
			goto done
		}
	}
done:
	url := strings.TrimSpace(rest[:end])
	return strings.TrimRight(url, ".,;!?)】）》」")
}

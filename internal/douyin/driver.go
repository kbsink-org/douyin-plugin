package douyinconv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/kbsink-org/kbsink/pkg/core"
)

var shareURLPattern = regexp.MustCompile(`https?://[^\s]+`)
var douyinVideoIDPattern = regexp.MustCompile(`(?i)(?:^|/)video/([^/?#]+)`)

// Driver fetches Douyin page HTML from a share link or share text.
type Driver struct {
	client        *http.Client
	videoPageBase string
}

func NewHTTPDriver(client *http.Client) *Driver {
	if client == nil {
		client = http.DefaultClient
	}
	return &Driver{
		client:        client,
		videoPageBase: "https://www.iesdouyin.com/share/video",
	}
}

func (d *Driver) Fetch(ctx context.Context, rawInput string) (*core.FetchResult, error) {
	shareURL := strings.TrimSpace(extractShareURL(rawInput))
	if shareURL == "" {
		return nil, core.NewCodedError(core.ErrCodeInvalidArgument, "douyin share link is required", nil)
	}

	initial, err := d.fetchHTML(ctx, shareURL)
	if err != nil {
		return nil, err
	}

	videoID := extractDouyinVideoID(initial.URL)
	if videoID == "" {
		return initial, nil
	}

	videoPageURL := strings.TrimRight(d.videoPageBase, "/") + "/" + videoID
	resolved, err := d.fetchHTML(ctx, videoPageURL)
	if err != nil {
		return initial, nil
	}
	return resolved, nil
}

func (d *Driver) fetchHTML(ctx context.Context, rawURL string) (*core.FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, core.NewCodedError(core.ErrCodeDriverBuildRequest, "build request", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1")
	req.Header.Set("Referer", "https://www.douyin.com/")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, core.NewCodedError(core.ErrCodeDriverRequestFailed, "execute request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, core.NewCodedError(
			core.ErrCodeDriverUnexpectedHTTP,
			fmt.Sprintf("unexpected status: %s", resp.Status),
			nil,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, core.NewCodedError(core.ErrCodeDriverReadBodyFailed, "read response body", err)
	}

	return &core.FetchResult{
		URL:  resp.Request.URL.String(),
		HTML: string(body),
	}, nil
}

func extractShareURL(rawInput string) string {
	rawInput = strings.TrimSpace(rawInput)
	if rawInput == "" {
		return ""
	}
	match := shareURLPattern.FindString(rawInput)
	return strings.TrimSpace(strings.TrimRight(match, ".,;!?)】）》」\"'"))
}

func extractDouyinVideoID(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	if qID := strings.TrimSpace(u.Query().Get("video_id")); qID != "" {
		return qID
	}

	pathValue := path.Clean(u.Path)
	match := douyinVideoIDPattern.FindStringSubmatch(pathValue)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

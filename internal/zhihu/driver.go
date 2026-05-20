package zhihuconv

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/kbsink-org/kbsink-plugins/internal/parseutil"
	"github.com/kbsink-org/kbsink/pkg/core"
)

const (
	zhihuUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	zhihuReferer   = "https://www.zhihu.com/"
)

// Driver fetches Zhihu zhuanlan article HTML.
type Driver struct {
	client core.HTTPClient
}

func NewHTTPDriver(client core.HTTPClient) *Driver {
	if client == nil {
		client = http.DefaultClient
	}
	return &Driver{client: client}
}

func (d *Driver) Fetch(ctx context.Context, rawInput string) (*core.FetchResult, error) {
	pageURL := parseutil.ExtractShareURL(rawInput)
	if pageURL == "" {
		pageURL = strings.TrimSpace(rawInput)
	}
	if pageURL == "" || !strings.Contains(strings.ToLower(pageURL), "zhihu.com") {
		return nil, core.NewCodedError(core.ErrCodeInvalidArgument, "zhihu article link is required", nil)
	}
	return d.fetchHTML(ctx, pageURL)
}

func (d *Driver) fetchHTML(ctx context.Context, rawURL string) (*core.FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, core.NewCodedError(core.ErrCodeDriverBuildRequest, "build request", err)
	}
	req.Header.Set("User-Agent", zhihuUserAgent)
	req.Header.Set("Referer", zhihuReferer)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, core.NewCodedError(core.ErrCodeDriverRequestFailed, "execute request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, core.NewCodedError(
			core.ErrCodeDriverUnexpectedHTTP,
			"unexpected status: "+resp.Status,
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

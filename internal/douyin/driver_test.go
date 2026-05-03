package douyinconv

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kbsink-org/kbsink/pkg/core"
)

func TestDriverFetch_ErrorCodes(t *testing.T) {
	d := NewHTTPDriver(nil)

	_, err := d.Fetch(context.Background(), "")
	if got := core.ErrorCodeOf(err); got != core.ErrCodeInvalidArgument {
		t.Fatalf("expected %s, got %s (err=%v)", core.ErrCodeInvalidArgument, got, err)
	}

	_, err = d.Fetch(context.Background(), "://bad-url")
	if got := core.ErrorCodeOf(err); got != core.ErrCodeInvalidArgument {
		t.Fatalf("expected %s, got %s (err=%v)", core.ErrCodeInvalidArgument, got, err)
	}
}

func TestDriverFetch_FromShareText(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/jump", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/video/123", http.StatusFound)
	})
	mux.HandleFunc("/video/123", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html><title>douyin</title></html>"))
	})
	mux.HandleFunc("/share/video/123", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html><title>douyin normalized</title></html>"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	shareText := "7.92 复制打开抖音看看 " + ts.URL + "/jump"
	d := NewHTTPDriver(ts.Client())
	d.videoPageBase = ts.URL + "/share/video"
	got, err := d.Fetch(context.Background(), shareText)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	if got == nil || got.HTML == "" {
		t.Fatalf("expected non-empty html")
	}
	if got.URL != ts.URL+"/share/video/123" {
		t.Fatalf("expected normalized video page url, got %q", got.URL)
	}
}

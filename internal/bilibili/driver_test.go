package bilibiliconv

import (
	"testing"
)

func TestNormalizeBilibiliVideoURL(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{
			"https://www.bilibili.com/video/BV1X1Lj66ENe/?share_source=copy_web",
			"https://www.bilibili.com/video/BV1X1Lj66ENe/",
		},
		{"https://b23.tv/BV1X1Lj66ENe", "https://www.bilibili.com/video/BV1X1Lj66ENe/"},
		{"text https://www.bilibili.com/video/bv1abc/ end", "https://www.bilibili.com/video/BV1abc/"},
		{"https://example.com/", ""},
	}
	for _, tt := range tests {
		got := normalizeBilibiliVideoURL(tt.in)
		if got != tt.want {
			t.Errorf("normalizeBilibiliVideoURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

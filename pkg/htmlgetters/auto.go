package htmlgetters

import (
	"context"
	"io"
	neturl "net/url"
	"strings"
)

// Get 获取 HTML
func Get(ctx context.Context, target string) (io.ReadCloser, *neturl.URL, error) {
	if target == "" {
		return Stdin()
	}
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return HTTP(ctx, target)
	}
	return File(target)
}

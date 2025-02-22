package htmlgetters

import (
	"context"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
)

// HTTP 通过 HTTP 获取网页内容
func HTTP(ctx context.Context, url string) (r io.ReadCloser, parsedURL *neturl.URL, err error) {
	// 解析 URL
	parsedURL, err = neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse url %q error: %w", url, err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("make request error: %w", err)
		return
	}

	// 请求
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		bodyRaw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		err = fmt.Errorf("received unexpected status code %d (!=200), body: %s", resp.StatusCode, string(bodyRaw))
		return
	}

	r = resp.Body
	return
}

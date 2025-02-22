package htmlgetters

import (
	"io"
	neturl "net/url"
	"os"
)

// Stdin 通过读标准输入获取 HTML
func Stdin() (io.ReadCloser, *neturl.URL, error) {
	return os.Stdin, &neturl.URL{}, nil
}

// File 通过读文件获取 HTML
func File(path string) (io.ReadCloser, *neturl.URL, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, &neturl.URL{Path: path}, nil
}

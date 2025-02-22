package articlewriters

import (
	"io"
	"strings"

	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
)

// NewHTML 创建以 HTML 格式写入文章的 Writer 的实现
func NewHTML(w io.Writer) Writer {
	return &htmlWriter{w: w}
}

// htmlWriter 以 HTML 格式写入文章的 Writer 的实现
type htmlWriter struct {
	w io.Writer
}

var _ Writer = (*htmlWriter)(nil)

// Write 写入文章
func (w *htmlWriter) Write(article readability.Article) error {
	cur := article.Node
	out := false
	indent := 0
	for cur != nil {
		switch cur.Type {
		case html.TextNode:
			if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+cur.Data+"\n"); err != nil {
				return err
			}
		default:
		}

		switch {
		case cur.FirstChild != nil && !out:
			// 往里走
			if cur.Type == html.ElementNode {
				attr := " "
				for _, a := range cur.Attr {
					attr += a.Key + "=\"" + a.Val + "\" "
				}
				attr = strings.TrimRight(attr, " ")
				if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+"<"+cur.Data+attr+">\n"); err != nil {
					return err
				}
			}
			cur = cur.FirstChild
			indent++
			out = false
		case cur.NextSibling != nil:
			// 下一个节点
			if cur.Type == html.ElementNode {
				if out {
					if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+"</"+cur.Data+">\n"); err != nil {
						return err
					}
				} else {
					attr := " "
					for _, a := range cur.Attr {
						attr += a.Key + "=\"" + a.Val + "\" "
					}
					attr = strings.TrimRight(attr, " ")
					if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+"<"+cur.Data+attr+"/>\n"); err != nil {
						return err
					}
				}
			}
			cur = cur.NextSibling
			out = false
		case cur.Parent != nil:
			// 返回上一级
			if cur.Type == html.ElementNode {
				if out {
					if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+"</"+cur.Data+">\n"); err != nil {
						return err
					}
				} else {
					attr := " "
					for _, a := range cur.Attr {
						attr += a.Key + "=\"" + a.Val + "\" "
					}
					attr = strings.TrimRight(attr, " ")
					if _, err := io.WriteString(w.w, strings.Repeat("  ", indent)+"<"+cur.Data+attr+"/>\n"); err != nil {
						return err
					}
				}
			}
			cur = cur.Parent
			indent--
			if indent < 0 {
				indent = 0
			}
			out = true
		default:
			// 结束
			cur = nil
		}
	}

	return nil
}

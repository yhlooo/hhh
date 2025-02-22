package articlewriters

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// NewMarkdown 创建以 Markdown 形式写入文章的 Writer 的实现
func NewMarkdown(w io.Writer) Writer {
	return &markdownWriter{w: w}
}

// markdownWriter 以 Markdown 形式写入文章的 Writer 的实现
type markdownWriter struct {
	w io.Writer
}

var _ Writer = (*markdownWriter)(nil)

// Write 写入文章
func (w *markdownWriter) Write(article readability.Article) error {
	if err := w.writeString(fmt.Sprintf("# %s\n\n", article.Title), ""); err != nil {
		return err
	}
	if err := w.writeNode(article.Node, &NodeOptions{}); err != nil {
		return err
	}
	if err := w.writeString("\n", ""); err != nil {
		return err
	}
	return nil
}

type NodeOptions struct {
	OrderedListSeq int
	Indent         string
}

// AppendIndent 返回追加了缩进的选项
func (opts *NodeOptions) AppendIndent(indent string) *NodeOptions {
	newOpts := *opts
	newOpts.Indent += indent
	return &newOpts
}

// WithOrderedListSeq 返回带指定 OrderedListSeq 的选项
func (opts *NodeOptions) WithOrderedListSeq(seq int) *NodeOptions {
	newOpts := *opts
	newOpts.OrderedListSeq = seq
	return &newOpts
}

// writeNode 写节点
func (w *markdownWriter) writeNode(node *html.Node, opts *NodeOptions) error {
	// 确定段落分隔符
	paragraphSeparator := ""
	if node.Type == html.ElementNode && node.NextSibling != nil {
		switch node.DataAtom {
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6,
			atom.P, atom.Div, atom.Figure, atom.Figcaption, atom.Section, atom.Ul, atom.Ol:
			paragraphSeparator = "\n\n"
		case atom.Li:
			paragraphSeparator = "\n"
		}
	}

	switch node.Type {
	case html.TextNode:
		return w.writeString(node.Data, opts.Indent)
	case html.ElementNode:
		switch node.DataAtom {
		case atom.H1:
			return w.writeChildren(node, "# ", paragraphSeparator, opts)
		case atom.H2:
			return w.writeChildren(node, "## ", paragraphSeparator, opts)
		case atom.H3:
			return w.writeChildren(node, "### ", paragraphSeparator, opts)
		case atom.H4:
			return w.writeChildren(node, "#### ", paragraphSeparator, opts)
		case atom.H5:
			return w.writeChildren(node, "##### ", paragraphSeparator, opts)
		case atom.H6:
			return w.writeChildren(node, "###### ", paragraphSeparator, opts)
		case atom.Figcaption:
			return w.writeChildren(node, "> ", paragraphSeparator, opts.AppendIndent("> "))
		case atom.Ol:
			return w.writeChildren(node, "", paragraphSeparator, opts.WithOrderedListSeq(1))
		case atom.Li:
			var err error
			if opts.OrderedListSeq == 0 {
				// 无序列表
				err = w.writeChildren(node, "- ", paragraphSeparator, opts.AppendIndent("  "))
			} else {
				// 有序列表
				err = w.writeChildren(
					node,
					strconv.Itoa(opts.OrderedListSeq)+". ",
					paragraphSeparator,
					opts.AppendIndent("   "),
				)
				opts.OrderedListSeq++
			}
			if err != nil {
				return err
			}
		case atom.B:
			return w.writeChildren(node, " **", "** ", opts)
		case atom.A:
			href := ""
			for _, attr := range node.Attr {
				switch attr.Key {
				case "href":
					href = attr.Val
				}
			}
			return w.writeChildren(node, " [", fmt.Sprintf("](%s) ", href), opts)
		case atom.Img, atom.Image:
			return w.writeImage(node)
		default:
			return w.writeChildren(node, "", paragraphSeparator, opts)
		}
	default:
	}

	return nil
}

// writeChildren 写子节点
func (w *markdownWriter) writeChildren(node *html.Node, prefix, suffix string, opts *NodeOptions) error {
	if prefix != "" {
		if err := w.writeString(prefix, ""); err != nil {
			return err
		}
	}
	for child := range node.ChildNodes() {
		if err := w.writeNode(child, opts); err != nil {
			return err
		}
	}
	if suffix != "" {
		if err := w.writeString(suffix, ""); err != nil {
			return err
		}
	}
	return nil
}

// writeString 写字符串
func (w *markdownWriter) writeString(text, indent string) error {
	if indent != "" {
		text = strings.ReplaceAll(text, "\n", "\n"+indent)
	}
	if _, err := io.WriteString(w.w, text); err != nil {
		return err
	}
	return nil
}

// writeImage 写图像节点
func (w *markdownWriter) writeImage(node *html.Node) error {
	src := ""
	alt := ""

	for _, attr := range node.Attr {
		switch attr.Key {
		case "src":
			src = attr.Val
		case "alt":
			alt = attr.Val
		}
	}

	paragraphSeparator := ""
	if node.NextSibling != nil {
		paragraphSeparator = "\n\n"
	}

	return w.writeString(fmt.Sprintf("![%s](%s)%s", alt, src, paragraphSeparator), "")
}

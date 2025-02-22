package articlewriters

import "github.com/go-shiori/go-readability"

// Writer 文章写入器
type Writer interface {
	// Write 写入文章
	Write(article readability.Article) error
}

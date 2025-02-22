package options

import "github.com/spf13/pflag"

// NewDefaultOptions 创建一个默认运行选项
func NewDefaultOptions() Options {
	return Options{
		Global: NewDefaultGlobalOptions(),
		Format: "markdown",
		Output: "",
	}
}

// Options 命令运行选项
type Options struct {
	// 全局选项
	Global GlobalOptions `json:"global,omitempty" yaml:"global,omitempty"`

	// 输出格式
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// 输出路径
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
}

// AddFlags 将选项绑定到命令行参数
func (opts *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Format, "format", "f", opts.Format, "Output format. One of: (markdown, html)")
	fs.StringVarP(&opts.Output, "output", "o", opts.Output, "Write to file instead of stdout")
}

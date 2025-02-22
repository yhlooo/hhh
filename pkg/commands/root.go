package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/go-shiori/go-readability"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yhlooo/hhh/pkg/articlewriters"
	"github.com/yhlooo/hhh/pkg/commands/options"
	"github.com/yhlooo/hhh/pkg/htmlgetters"
)

// NewCommand 创建命令
func NewCommand() *cobra.Command {
	opts := options.NewDefaultOptions()
	return NewCommandWithOptions(&opts)
}

// NewCommandWithOptions 基于选项创建命令
func NewCommandWithOptions(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hhh [flags] [url|file]",
		Short: "Human readable HTML Helper",
		Long: `Human readable HTML Helper

Retrieve HTML content via HTTP requests or local file reading, denoise and
extract the primary textual content from the page, then render the output
in human-readable format (e.g. Markdown).
`,
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 校验全局选项
			if err := opts.Global.Validate(); err != nil {
				return err
			}
			// 设置日志
			logger := setLogger(cmd, opts.Global.Verbosity)
			// 输出选项
			optsRaw, _ := json.Marshal(opts)
			logger.V(1).Info(fmt.Sprintf("command: %q, args: %q, options: %s", cmd.Name(), args, string(optsRaw)))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			targets := args
			if len(targets) == 0 {
				targets = []string{""}
			}

			// 准备输出
			w := os.Stdout
			if opts.Output != "" {
				var err error
				w, err = os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return fmt.Errorf("open %q error: %w", opts.Output, err)
				}
			}
			defer func() { _ = w.Close() }()
			var articleWriter articlewriters.Writer
			switch opts.Format {
			case "markdown":
				articleWriter = articlewriters.NewMarkdown(w)
			default:
				articleWriter = articlewriters.NewHTML(w)
			}

			var errs []error
			for _, target := range targets {
				// 获取 HTML
				r, parsedURL, err := htmlgetters.Get(ctx, target)
				if err != nil {
					logger.Error(err, fmt.Sprintf("get %q error", target))
					errs = append(errs, fmt.Errorf("get %q error: %w", target, err))
					continue
				}

				// 降噪
				article, err := readability.FromReader(r, parsedURL)
				if err != nil {
					logger.Error(err, fmt.Sprintf("make %q readable error", target))
					errs = append(errs, fmt.Errorf("make %q readable error: %w", target, err))
					continue
				}
				_ = r.Close()

				// 输出
				if err := articleWriter.Write(article); err != nil {
					logger.Error(err, fmt.Sprintf("write %q to output error", target))
					errs = append(errs, fmt.Errorf("write %q to output error: %w", target, err))
					continue
				}
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			return nil
		},
	}

	// 将选项绑定到命令行
	opts.Global.AddFlags(cmd.PersistentFlags())
	opts.AddFlags(cmd.Flags())

	return cmd
}

// setLogger 设置命令日志，并返回 logr.Logger
func setLogger(cmd *cobra.Command, verbosity uint32) logr.Logger {
	// 设置日志级别
	logrusLogger := logrus.New()
	switch verbosity {
	case 1:
		logrusLogger.SetLevel(logrus.DebugLevel)
	case 2:
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}
	// 将 logger 注入上下文
	logger := logrusr.New(logrusLogger)
	cmd.SetContext(logr.NewContext(cmd.Context(), logger))

	return logger
}

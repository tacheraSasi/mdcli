package renderer

import (
	"bytes"
	"os"

	"github.com/MichaelMure/go-term-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type RenderOptions struct {
	Input    string
	Autolink bool
}

func Render(opts RenderOptions) (string, error) {
	var md goldmark.Markdown
	if opts.Autolink {
		md = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
				extension.NewLinkify(
					extension.WithLinkifyAllowedProtocols([][]byte{
						[]byte("http:"),
						[]byte("https:"),
					}),
				),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)
	} else {
		md = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(opts.Input), &buf); err != nil {
		return "", err
	}

	result := markdown.Render(buf.String(), 80, 6)

	return string(result), nil
}

func ReadFile(file string) (string, error) {
	// Reading file
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
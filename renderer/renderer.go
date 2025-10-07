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
	Input        string
	Autolink     bool
	Theme        string
	Width        int
	OutputFormat string
}

func Render(opts RenderOptions) (string, error) {
	// Set defaults
	if opts.Theme == "" {
		opts.Theme = "dracula"
	}
	if opts.Width == 0 {
		opts.Width = 80
	}
	if opts.OutputFormat == "" {
		opts.OutputFormat = "terminal"
	}

	var md goldmark.Markdown
	extensions := []goldmark.Extender{
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle(opts.Theme),
		),
	}

	if opts.Autolink {
		extensions = append(extensions, extension.NewLinkify(
			extension.WithLinkifyAllowedProtocols([][]byte{
				[]byte("http:"),
				[]byte("https:"),
			}),
		))
	}

	md = goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(opts.Input), &buf); err != nil {
		return "", err
	}

	switch opts.OutputFormat {
	case "html":
		return buf.String(), nil
	case "text", "plain":
		return stripHTML(buf.String()), nil
	case "pdf":
		return renderToPDF(buf.String())
	default: // terminal
		result := markdown.Render(buf.String(), opts.Width, 6)
		return string(result), nil
	}
}

// stripHTML removes HTML tags from content for plain text output
func stripHTML(content string) string {
	// Simple HTML tag removal - you might want to use a proper HTML parser
	var result strings.Builder
	inTag := false
	
	for _, char := range content {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(char)
		}
	}
	
	return result.String()
}

// renderToPDF converts HTML to PDF (simplified implementation)
func renderToPDF(htmlContent string) (string, error) {
	// This is a placeholder for PDF rendering
	// In a real implementation, you'd use a library like wkhtmltopdf or similar
	return "PDF rendering not implemented in this version", nil
}

func ReadFile(file string) (string, error) {
	// Reading file
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
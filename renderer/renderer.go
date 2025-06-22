package renderer

import (
	"os"

	"github.com/charmbracelet/glamour"
)

func Render(input string, style string) (string, error) {
	// Rendering Markdown with glamour
	rendered, err := glamour.Render(input, style)
	if err != nil {
		return "", err
	}
	return rendered, nil
}

func ReadFile(file string) (string, error) {
	// Reading file
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

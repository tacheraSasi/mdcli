package renderer

import "github.com/charmbracelet/glamour"

func Render(input string) (string, error) {
	// Rendering Markdown with glamour
	rendered, err := glamour.Render(input, "dark") 
	if err != nil {
		return "", err
	}
	return rendered, nil
}
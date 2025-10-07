package themes

import (
	"fmt"
)

// Theme represents a color theme for terminal output
type Theme struct {
	Name        string
	Description string
	Colors      ThemeColors
}

// ThemeColors defines the color scheme
type ThemeColors struct {
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Text       string
	Link       string
	Code       string
	Header     string
}

// Available themes
var AvailableThemes = map[string]Theme{
	"dracula": {
		Name:        "Dracula",
		Description: "Dark theme with purple accents",
		Colors: ThemeColors{
			Primary:    "#bd93f9",
			Secondary:  "#6272a4",
			Accent:     "#50fa7b",
			Background: "#282a36",
			Text:       "#f8f8f2",
			Link:       "#8be9fd",
			Code:       "#ffb86c",
			Header:     "#ff79c6",
		},
	},
	"github": {
		Name:        "GitHub",
		Description: "Light theme inspired by GitHub",
		Colors: ThemeColors{
			Primary:    "#0366d6",
			Secondary:  "#586069",
			Accent:     "#28a745",
			Background: "#ffffff",
			Text:       "#24292e",
			Link:       "#0366d6",
			Code:       "#d73a49",
			Header:     "#24292e",
		},
	},
	"monokai": {
		Name:        "Monokai",
		Description: "Dark theme with vibrant colors",
		Colors: ThemeColors{
			Primary:    "#f92672",
			Secondary:  "#75715e",
			Accent:     "#a6e22e",
			Background: "#272822",
			Text:       "#f8f8f2",
			Link:       "#66d9ef",
			Code:       "#fd971f",
			Header:     "#f92672",
		},
	},
	"solarized": {
		Name:        "Solarized",
		Description: "Balanced light/dark theme",
		Colors: ThemeColors{
			Primary:    "#268bd2",
			Secondary:  "#93a1a1",
			Accent:     "#859900",
			Background: "#fdf6e3",
			Text:       "#657b83",
			Link:       "#268bd2",
			Code:       "#d33682",
			Header:     "#b58900",
		},
	},
	"nord": {
		Name:        "Nord",
		Description: "Arctic, north-bluish color palette",
		Colors: ThemeColors{
			Primary:    "#5e81ac",
			Secondary:  "#4c566a",
			Accent:     "#a3be8c",
			Background: "#2e3440",
			Text:       "#d8dee9",
			Link:       "#88c0d0",
			Code:       "#ebcb8b",
			Header:     "#81a1c1",
		},
	},
}

// GetTheme returns a theme by name
func GetTheme(name string) (Theme, error) {
	if theme, exists := AvailableThemes[name]; exists {
		return theme, nil
	}
	return Theme{}, fmt.Errorf("theme '%s' not found", name)
}

// ListThemes returns all available theme names
func ListThemes() []string {
	var themes []string
	for name := range AvailableThemes {
		themes = append(themes, name)
	}
	return themes
}

// GetSyntaxHighlightingStyle returns the appropriate syntax highlighting style for a theme
func GetSyntaxHighlightingStyle(themeName string) string {
	switch themeName {
	case "github":
		return "github"
	case "monokai":
		return "monokai"
	case "solarized":
		return "solarized-light"
	case "nord":
		return "nord"
	default:
		return "dracula"
	}
}
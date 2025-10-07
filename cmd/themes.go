package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/themes"
)

var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "List available themes",
	Long: `Display all available color themes for syntax highlighting and terminal output.
Each theme has its own color palette optimized for different preferences.`,
	Run: runThemes,
}

func init() {
	rootCmd.AddCommand(themesCmd)
}

func runThemes(cmd *cobra.Command, args []string) {
	fmt.Println(" Available Themes:")
	fmt.Println(strings.Repeat("=", 40))

	for name, theme := range themes.AvailableThemes {
		fmt.Printf("\n%s\n", strings.ToUpper(name))
		fmt.Printf("   Name: %s\n", theme.Name)
		fmt.Printf("   Description: %s\n", theme.Description)

		if verbose {
			fmt.Printf("   Colors:\n")
			fmt.Printf("     Primary: %s\n", theme.Colors.Primary)
			fmt.Printf("     Accent: %s\n", theme.Colors.Accent)
			fmt.Printf("     Background: %s\n", theme.Colors.Background)
			fmt.Printf("     Text: %s\n", theme.Colors.Text)
		}
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("Use --theme=<name> with render command to apply a theme")
	fmt.Println("Use --verbose to see color details")
}

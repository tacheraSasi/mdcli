package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/renderer"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive Markdown editor mode",
	Long: `Start an interactive mode where you can input Markdown content 
and see the rendered output immediately. Type 'exit' or 'quit' to leave 
interactive mode.`,
	Aliases: []string{"i", "repl"},
	Run:     runInteractive,
}

var (
	interactiveFormat string
	interactiveTheme  string
	interactiveWidth  int
)

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.Flags().StringVarP(&interactiveFormat, "format", "f", "terminal", "Output format")
	interactiveCmd.Flags().StringVarP(&interactiveTheme, "theme", "t", "dracula", "Syntax highlighting theme")
	interactiveCmd.Flags().IntVarP(&interactiveWidth, "width", "w", 80, "Terminal width")
}

func runInteractive(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Welcome to mdcli Interactive Mode!")
	fmt.Println("Enter Markdown content and press Ctrl+D (EOF) to render.")
	fmt.Println("Commands: 'exit', 'quit', 'help', 'clear'")
	fmt.Println(strings.Repeat("-", 50))

	scanner := bufio.NewScanner(os.Stdin)
	var buffer strings.Builder

	for {
		fmt.Print("mdcli> ")

		// Read multiple lines until EOF or empty line
		buffer.Reset()

		for scanner.Scan() {
			line := scanner.Text()

			// Handle special commands
			if buffer.Len() == 0 {
				switch strings.ToLower(strings.TrimSpace(line)) {
				case "exit", "quit":
					fmt.Println("👋 Goodbye!")
					return
				case "help":
					showInteractiveHelp()
					continue
				case "clear":
					fmt.Print("\033[2J\033[H") // Clear screen
					continue
				}
			}

			if line == "" && buffer.Len() > 0 {
				// Empty line signals end of input
				break
			}

			if buffer.Len() > 0 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		content := buffer.String()
		if content == "" {
			continue
		}

		// Render the content
		rendered, err := renderer.Render(renderer.RenderOptions{
			Input:        content,
			Autolink:     true,
			Theme:        interactiveTheme,
			Width:        interactiveWidth,
			OutputFormat: interactiveFormat,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error rendering: %v\n", err)
			continue
		}

		// Display the result
		fmt.Println("\n" + strings.Repeat("=", 40))
		fmt.Println("📄 Rendered Output:")
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println(rendered)
		fmt.Println(strings.Repeat("-", 40))
	}
}

func showInteractiveHelp() {
	help := `
Interactive Mode Help:
=====================

Commands:
  exit, quit    - Exit interactive mode
  help          - Show this help message
  clear         - Clear the terminal screen

Usage:
  - Type your Markdown content
  - Press Enter twice (empty line) to render
  - Use Ctrl+C to exit anytime

Supported Markdown:
  - Headers (# ## ###)
  - **Bold** and *italic* text
  - Links [text](url)
  - Code blocks and inline code
  - Lists and tables
  - And much more!

Current settings:
  Format: %s
  Theme:  %s
  Width:  %d

`
	fmt.Printf(help, interactiveFormat, interactiveTheme, interactiveWidth)
}

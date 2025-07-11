package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tacheraSasi/mdcli/renderer"
)

var (
	version    = "1.0.0"
	style      = flag.String("style", "dark", "Style to use for rendering")
	output     = flag.String("output", "", "Write rendered output to file (default: stdout)")
	listStyles = flag.Bool("list-styles", false, "List available styles and exit")
)

var availableStyles = []string{
	"dark",
	"light",
	"notty",
	"pink",
	"solarized-dark",
	"solarized-light",
	"dracula",
	"no-color",
	"auto",
}

func main() {
	usageStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	exampleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("36"))

	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, usageStyle.Render(fmt.Sprintf("Usage: %s [options] [file ...]", "mdcli")))
		fmt.Fprintln(os.Stdout, headerStyle.Render("Options:"))
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, headerStyle.Render("Examples:"))
		fmt.Fprintln(os.Stdout, exampleStyle.Render("  mdcli README.md"))
		fmt.Fprintln(os.Stdout, exampleStyle.Render("  mdcli --style=dracula notes.md"))
		fmt.Fprintln(os.Stdout, exampleStyle.Render("  mdcli --output=out.txt file1.md file2.md"))
		fmt.Fprintln(os.Stdout, exampleStyle.Render("  cat README.md | mdcli --style=light"))
		fmt.Fprintln(os.Stdout, exampleStyle.Render("  mdcli --list-styles"))
	}

	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if *listStyles {
		fmt.Println("Available styles:")
		for _, s := range availableStyles {
			fmt.Println("  -", s)
		}
		return
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var inputs []string
	var filenames []string

	if flag.NArg() == 0 {
		// Read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		inputs = append(inputs, string(content))
	} else {
		for i := 0; i < flag.NArg(); i++ {
			filename := flag.Arg(i)
			if !strings.HasSuffix(filename, ".md") {
				return fmt.Errorf("file must have a .md extension: %s", filename)
			}
			content, err := renderer.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", filename, err)
			}
			inputs = append(inputs, content)
			filenames = append(filenames, filename)
		}
	}

	var renderedAll []string
	for idx, input := range inputs {
		rendered, err := renderer.Render(input, *style)
		if err != nil {
			return fmt.Errorf("error rendering markdown for %s: %w", func() string {
				if len(filenames) > idx {
					return filenames[idx]
				} else {
					return "stdin"
				}
			}(), err)
		}
		renderedAll = append(renderedAll, rendered)
	}

	outputStr := strings.Join(renderedAll, "\n\n---\n\n")

	if *output != "" {
		err := os.WriteFile(*output, []byte(outputStr), 0644)
		if err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
	} else {
		fmt.Println(outputStr)
	}
	return nil
}

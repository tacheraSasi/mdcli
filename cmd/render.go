package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tacheraSasi/mdcli/renderer"
)

var renderCmd = &cobra.Command{
	Use:   "render [files...]",
	Short: "Render Markdown files to various formats",
	Long: `Render one or more Markdown files to the specified output format.
Supports terminal output (default), HTML, PDF, and plain text formats.
Can process multiple files and supports stdin input.`,
	Args: cobra.ArbitraryArgs,
	Run:  runRender,
}

var (
	outputFile   string
	outputFormat string
	theme        string
	width        int
	autolink     bool
	showProgress bool
)

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	renderCmd.Flags().StringVarP(&outputFormat, "format", "f", "terminal", "Output format (terminal, html, pdf, text)")
	renderCmd.Flags().StringVarP(&theme, "theme", "t", "", "Syntax highlighting theme")
	renderCmd.Flags().IntVarP(&width, "width", "w", 0, "Terminal width for formatting")
	renderCmd.Flags().BoolVar(&autolink, "autolink", true, "Enable automatic link detection")
	renderCmd.Flags().BoolVar(&showProgress, "progress", false, "Show progress bar")

	// Bind flags to viper
	viper.BindPFlag("output", renderCmd.Flags().Lookup("output"))
	viper.BindPFlag("format", renderCmd.Flags().Lookup("format"))
	viper.BindPFlag("theme", renderCmd.Flags().Lookup("theme"))
	viper.BindPFlag("width", renderCmd.Flags().Lookup("width"))
	viper.BindPFlag("autolink", renderCmd.Flags().Lookup("autolink"))
}

func runRender(cmd *cobra.Command, args []string) {
	// Get values from viper (config file or flags)
	if outputFile == "" {
		outputFile = viper.GetString("output")
	}
	if outputFormat == "terminal" {
		outputFormat = viper.GetString("output_format")
	}
	if theme == "" {
		theme = viper.GetString("theme")
	}
	if width == 0 {
		width = viper.GetInt("width")
	}
	if !cmd.Flags().Changed("autolink") {
		autolink = viper.GetBool("autolink")
	}

	var inputs []string
	var filenames []string

	if len(args) == 0 {
		// Read from stdin
		if verbose {
			fmt.Fprintln(os.Stderr, "Reading from stdin...")
		}
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		inputs = append(inputs, string(content))
		filenames = append(filenames, "stdin")
	} else {
		// Validate and read files
		for _, filename := range args {
			if !strings.HasSuffix(strings.ToLower(filename), ".md") &&
				!strings.HasSuffix(strings.ToLower(filename), ".markdown") {
				fmt.Fprintf(os.Stderr, "Warning: %s doesn't appear to be a Markdown file\n", filename)
			}

			content, err := renderer.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
				os.Exit(1)
			}
			inputs = append(inputs, content)
			filenames = append(filenames, filename)
		}
	}

	// Setup progress bar if requested
	var bar *progressbar.ProgressBar
	if showProgress && len(inputs) > 1 {
		bar = progressbar.NewOptions(len(inputs),
			progressbar.OptionSetDescription("Processing files..."),
			progressbar.OptionSetWidth(15),
			progressbar.OptionShowCount(),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	var renderedAll []string
	for idx, input := range inputs {
		if verbose {
			fmt.Fprintf(os.Stderr, "Processing: %s\n", filenames[idx])
		}

		rendered, err := renderer.Render(renderer.RenderOptions{
			Input:        input,
			Autolink:     autolink,
			Theme:        theme,
			Width:        width,
			OutputFormat: outputFormat,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering %s: %v\n", filenames[idx], err)
			os.Exit(1)
		}
		renderedAll = append(renderedAll, rendered)

		if bar != nil {
			bar.Add(1)
			time.Sleep(10 * time.Millisecond) // Small delay for better UX
		}
	}

	// Join all rendered content
	var outputStr string
	if len(renderedAll) == 1 {
		outputStr = renderedAll[0]
	} else {
		switch outputFormat {
		case "html":
			outputStr = strings.Join(renderedAll, "\n<hr>\n")
		case "pdf":
			outputStr = strings.Join(renderedAll, "\n\n---\n\n")
		default:
			outputStr = strings.Join(renderedAll, "\n\n---\n\n")
		}
	}

	// Output the result
	if outputFile != "" {
		// Ensure output directory exists
		if dir := filepath.Dir(outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		err := os.WriteFile(outputFile, []byte(outputStr), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Output written to: %s\n", outputFile)
		}
	} else {
		fmt.Print(outputStr)
	}
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/renderer"
)

var watchCmd = &cobra.Command{
	Use:   "watch [file]",
	Short: "Watch a Markdown file for changes and auto-regenerate output",
	Long: `Watch one or more Markdown files for changes and automatically 
regenerate the output when files are modified. Perfect for live preview 
during document editing.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runWatch,
}

var (
	watchOutput string
	watchFormat string
	watchTheme  string
	watchWidth  int
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&watchOutput, "output", "o", "", "Output file path")
	watchCmd.Flags().StringVarP(&watchFormat, "format", "f", "terminal", "Output format")
	watchCmd.Flags().StringVarP(&watchTheme, "theme", "t", "dracula", "Syntax highlighting theme")
	watchCmd.Flags().IntVarP(&watchWidth, "width", "w", 80, "Terminal width")
}

func runWatch(cmd *cobra.Command, args []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file watcher: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	// Add files to watcher
	filesToWatch := make(map[string]bool)
	for _, file := range args {
		absPath, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting absolute path for %s: %v\n", file, err)
			continue
		}

		err = watcher.Add(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error watching file %s: %v\n", file, err)
			continue
		}

		filesToWatch[absPath] = true
		if verbose {
			fmt.Fprintf(os.Stderr, "Watching: %s\n", absPath)
		}
	}

	// Initial render
	renderFiles(args)

	fmt.Println("👀 Watching for changes... Press Ctrl+C to stop.")

	// Watch for file changes
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				if filesToWatch[event.Name] {
					if verbose {
						fmt.Fprintf(os.Stderr, "File modified: %s\n", event.Name)
					}

					// Add a small delay to avoid multiple rapid updates
					time.Sleep(100 * time.Millisecond)
					renderFiles(args)
					fmt.Println("✅ Updated!")
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
		}
	}
}

func renderFiles(files []string) {
	var inputs []string
	var filenames []string

	for _, filename := range files {
		if !strings.HasSuffix(strings.ToLower(filename), ".md") &&
			!strings.HasSuffix(strings.ToLower(filename), ".markdown") {
			continue
		}

		content, err := renderer.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
			continue
		}
		inputs = append(inputs, content)
		filenames = append(filenames, filename)
	}

	var renderedAll []string
	for idx, input := range inputs {
		rendered, err := renderer.Render(renderer.RenderOptions{
			Input:        input,
			Autolink:     true,
			Theme:        watchTheme,
			Width:        watchWidth,
			OutputFormat: watchFormat,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering %s: %v\n", filenames[idx], err)
			continue
		}
		renderedAll = append(renderedAll, rendered)
	}

	outputStr := strings.Join(renderedAll, "\n\n---\n\n")

	if watchOutput != "" {
		err := os.WriteFile(watchOutput, []byte(outputStr), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
		}
	} else {
		// Clear screen before showing new output
		fmt.Print("\033[2J\033[H")
		fmt.Print(outputStr)
	}
}

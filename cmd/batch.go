package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/renderer"
)

var batchCmd = &cobra.Command{
	Use:   "batch [directory]",
	Short: "Process all Markdown files in a directory",
	Long: `Batch process all Markdown files in a directory and its subdirectories.
Supports concurrent processing for better performance and various output formats.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runBatch,
}

var (
	batchOutput     string
	batchFormat     string
	batchTheme      string
	batchWidth      int
	batchRecursive  bool
	batchConcurrent int
	batchExtension  string
)

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&batchOutput, "output", "o", "", "Output directory")
	batchCmd.Flags().StringVarP(&batchFormat, "format", "f", "html", "Output format")
	batchCmd.Flags().StringVarP(&batchTheme, "theme", "t", "dracula", "Syntax highlighting theme")
	batchCmd.Flags().IntVarP(&batchWidth, "width", "w", 80, "Terminal width")
	batchCmd.Flags().BoolVarP(&batchRecursive, "recursive", "r", false, "Process subdirectories recursively")
	batchCmd.Flags().IntVarP(&batchConcurrent, "concurrent", "c", 4, "Number of concurrent workers")
	batchCmd.Flags().StringVarP(&batchExtension, "ext", "e", ".html", "Output file extension")
}

type BatchJob struct {
	InputFile  string
	OutputFile string
	Content    string
}

func runBatch(cmd *cobra.Command, args []string) {
	inputDir := args[0]

	// Find all markdown files
	var markdownFiles []string
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !batchRecursive && path != inputDir {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), ".md") ||
			strings.HasSuffix(strings.ToLower(path), ".markdown") {
			markdownFiles = append(markdownFiles, path)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(markdownFiles) == 0 {
		fmt.Println("No Markdown files found in the specified directory.")
		return
	}

	fmt.Printf("Found %d Markdown files\n", len(markdownFiles))

	// Prepare output directory
	outputDir := batchOutput
	if outputDir == "" {
		outputDir = filepath.Join(inputDir, "output")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Create job queue
	jobs := make(chan BatchJob, len(markdownFiles))
	results := make(chan error, len(markdownFiles))

	// Progress bar
	bar := progressbar.NewOptions(len(markdownFiles),
		progressbar.OptionSetDescription("Processing files..."),
		progressbar.OptionSetWidth(15),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
	)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < batchConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				err := processBatchJob(job)
				results <- err
				bar.Add(1)
			}
		}()
	}

	// Queue jobs
	go func() {
		defer close(jobs)
		for _, file := range markdownFiles {
			relPath, _ := filepath.Rel(inputDir, file)
			outputFile := filepath.Join(outputDir,
				strings.TrimSuffix(relPath, filepath.Ext(relPath))+batchExtension)

			// Ensure output subdirectory exists
			if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output subdirectory: %v\n", err)
				continue
			}

			content, err := renderer.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
				continue
			}

			jobs <- BatchJob{
				InputFile:  file,
				OutputFile: outputFile,
				Content:    content,
			}
		}
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var errorCount int
	for err := range results {
		if err != nil {
			errorCount++
			if verbose {
				fmt.Fprintf(os.Stderr, "Processing error: %v\n", err)
			}
		}
	}

	bar.Finish()
	fmt.Printf("\n✅ Processed %d files", len(markdownFiles)-errorCount)
	if errorCount > 0 {
		fmt.Printf(" (%d errors)", errorCount)
	}
	fmt.Printf("\n📁 Output directory: %s\n", outputDir)
}

func processBatchJob(job BatchJob) error {
	rendered, err := renderer.Render(renderer.RenderOptions{
		Input:        job.Content,
		Autolink:     true,
		Theme:        batchTheme,
		Width:        batchWidth,
		OutputFormat: batchFormat,
	})
	if err != nil {
		return fmt.Errorf("error rendering %s: %w", job.InputFile, err)
	}

	err = os.WriteFile(job.OutputFile, []byte(rendered), 0644)
	if err != nil {
		return fmt.Errorf("error writing %s: %w", job.OutputFile, err)
	}

	return nil
}

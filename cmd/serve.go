package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/a-h/templ"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/renderer"
	"github.com/tacheraSasi/mdcli/themes"
	views "github.com/tacheraSasi/mdcli/ui"
)

var serveCmd = &cobra.Command{
	Use:   "serve [file]",
	Short: "Start a live preview server for Markdown files",
	Long: `Start a web server that provides live preview of Markdown files.
The preview updates automatically when files are modified.
Perfect for real-time document editing and review.`,
	Args: cobra.ExactArgs(1),
	Run:  runServe,
}

var (
	servePort   int
	serveTheme  string
	serveWidth  int
	serveBind   string
	serveReload bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to serve on")
	serveCmd.Flags().StringVarP(&serveTheme, "theme", "t", "github", "Theme for HTML output")
	serveCmd.Flags().IntVarP(&serveWidth, "width", "w", 80, "Content width")
	serveCmd.Flags().StringVarP(&serveBind, "bind", "b", "localhost", "Bind address")
	serveCmd.Flags().BoolVar(&serveReload, "auto-reload", true, "Enable auto-reload on file changes")
}

type PreviewData = views.ServeData

var (
	currentFile   string
	lastModTime   time.Time
	cachedContent string
)

func runServe(cmd *cobra.Command, args []string) {
	currentFile = args[0]

	// Check if file exists
	if _, err := os.Stat(currentFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File not found: %s\n", currentFile)
		os.Exit(1)
	}

	// Parse template
	var err error
	previewTemplate, err = template.New("preview").Parse(htmlTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Get theme
	theme, err := themes.GetTheme(serveTheme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Theme error: %v\n", err)
		os.Exit(1)
	}

	// Initial render
	if err := renderCurrentFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Initial render error: %v\n", err)
		os.Exit(1)
	}

	// Setup file watcher if auto-reload is enabled
	if serveReload {
		go startFileWatcher()
	}

	// Setup HTTP handlers
	data := PreviewData{
		Title:      filepath.Base(currentFile),
		Content:    template.HTML(cachedContent),
		Theme:      theme,
		AutoReload: serveReload,
	}
	http.Handle("/", templ.Handler(views.Serve(data)))

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"lastModified": %d}`, lastModTime.Unix())
	})

	addr := fmt.Sprintf("%s:%d", serveBind, servePort)
	fmt.Printf("🚀 Starting live preview server...\n")
	fmt.Printf("📄 File: %s\n", currentFile)
	fmt.Printf("🌐 URL: http://%s\n", addr)
	fmt.Printf("🎨 Theme: %s\n", serveTheme)
	if serveReload {
		fmt.Printf("🔄 Auto-reload: enabled\n")
	}
	fmt.Printf("Press Ctrl+C to stop\n\n")

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func renderCurrentFile() error {
	stat, err := os.Stat(currentFile)
	if err != nil {
		return err
	}

	lastModTime = stat.ModTime()

	content, err := renderer.ReadFile(currentFile)
	if err != nil {
		return err
	}

	rendered, err := renderer.Render(renderer.RenderOptions{
		Input:        content,
		Autolink:     true,
		Theme:        serveTheme,
		Width:        serveWidth,
		OutputFormat: "html",
	})
	if err != nil {
		return err
	}

	cachedContent = rendered
	return nil
}

func startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	absPath, _ := filepath.Abs(currentFile)
	err = watcher.Add(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error watching file: %v\n", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				time.Sleep(100 * time.Millisecond) // Debounce
				if err := renderCurrentFile(); err != nil {
					fmt.Fprintf(os.Stderr, "Render error: %v\n", err)
				} else if verbose {
					fmt.Printf("📝 File updated: %s\n", time.Now().Format("15:04:05"))
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

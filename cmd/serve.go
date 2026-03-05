package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/tacheraSasi/mdcli/renderer"
	views "github.com/tacheraSasi/mdcli/ui"
)

// AssetsFS is set from main to provide embedded assets.
var AssetsFS embed.FS

var serveCmd = &cobra.Command{
	Use:   "serve [file-or-directory]",
	Short: "Start a live preview server for Markdown files",
	Long: `Start a web server that provides live preview of Markdown files.
Supports both single file and directory serving (including nested folders).
The preview updates automatically when files are modified.
Perfect for real-time document editing and review.

Examples:
  mdcli serve README.md          # Serve a single file
  mdcli serve .                  # Serve all .md files in current directory
  mdcli serve docs/              # Serve all .md files in docs/ recursively`,
	Args: cobra.MaximumNArgs(1),
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

// --- Single-file mode state ---
var (
	currentFile   string
	lastModTime   time.Time
	cachedContent string
)

// --- Directory mode state ---
var (
	isDirectoryMode bool
	baseDir         string
	fileCache       map[string]*CachedFile
	fileCacheMu     sync.RWMutex
	fileTree        []views.FileEntry
	globalModTime   time.Time
)

// CachedFile stores the rendered HTML and modification time for a single file.
type CachedFile struct {
	Content string
	ModTime time.Time
}

// Directories to skip when scanning.
var skipDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	".git":         true,
	"batch_output": true,
	"build":        true,
	"bin":          true,
	".github":      true,
	"__pycache__":  true,
}

func runServe(cmd *cobra.Command, args []string) {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	stat, err := os.Stat(target)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Path not found: %s\n", target)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if stat.IsDir() {
		isDirectoryMode = true
		runServeDirectory(cmd, target)
	} else {
		isDirectoryMode = false
		runServeSingleFile(cmd, target)
	}
}

// ====================================================================
// SINGLE-FILE MODE
// ====================================================================

func runServeSingleFile(cmd *cobra.Command, file string) {
	currentFile = file

	if err := renderCurrentFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Initial render error: %v\n", err)
		os.Exit(1)
	}

	if serveReload {
		go startFileWatcher()
	}

	mux := http.NewServeMux()

	assetsSubFS, _ := fs.Sub(AssetsFS, "assets")
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsSubFS))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := views.ServeData{
			Title:      filepath.Base(currentFile),
			Content:    cachedContent,
			ThemeName:  serveTheme,
			AutoReload: serveReload,
		}
		templ.Handler(views.ServePage(data)).ServeHTTP(w, r)
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
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

	if err := http.ListenAndServe(addr, mux); err != nil {
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

// ====================================================================
// DIRECTORY MODE
// ====================================================================

func runServeDirectory(cmd *cobra.Command, dir string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}
	baseDir = absDir
	fileCache = make(map[string]*CachedFile)

	// Initial scan and render all markdown files
	if err := scanAndRenderDirectory(); err != nil {
		fmt.Fprintf(os.Stderr, "Initial scan error: %v\n", err)
		os.Exit(1)
	}

	if serveReload {
		go startDirectoryWatcher()
	}

	mux := http.NewServeMux()

	assetsSubFS, _ := fs.Sub(AssetsFS, "assets")
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsSubFS))))

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fileCacheMu.RLock()
		modTime := globalModTime
		fileCacheMu.RUnlock()
		fmt.Fprintf(w, `{"lastModified": %d}`, modTime.Unix())
	})

	mux.HandleFunc("/", handleDirectoryRequest)

	addr := fmt.Sprintf("%s:%d", serveBind, servePort)
	fileCacheMu.RLock()
	fileCount := len(fileCache)
	fileCacheMu.RUnlock()

	fmt.Printf("🚀 Starting live preview server (directory mode)...\n")
	fmt.Printf("📁 Directory: %s\n", baseDir)
	fmt.Printf("📄 Files: %d markdown files found\n", fileCount)
	fmt.Printf("🌐 URL: http://%s\n", addr)
	fmt.Printf("🎨 Theme: %s\n", serveTheme)
	if serveReload {
		fmt.Printf("🔄 Auto-reload: enabled\n")
	}
	fmt.Printf("Press Ctrl+C to stop\n\n")

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// handleDirectoryRequest routes incoming requests to the correct markdown
// file or shows a directory listing.
func handleDirectoryRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")

	fileCacheMu.RLock()
	tree := fileTree
	cache := make(map[string]*CachedFile, len(fileCache))
	for k, v := range fileCache {
		cache[k] = v
	}
	fileCacheMu.RUnlock()

	var content string
	var title string
	var currentPath string
	found := false

	// 1. Exact match (e.g., /README.md)
	if cached, ok := cache[urlPath]; ok {
		content = cached.Content
		title = filepath.Base(urlPath)
		currentPath = urlPath
		found = true
	}

	// 2. Try adding .md
	if !found {
		mdPath := urlPath + ".md"
		if cached, ok := cache[mdPath]; ok {
			content = cached.Content
			title = filepath.Base(mdPath)
			currentPath = mdPath
			found = true
		}
	}

	// 3. Try adding .markdown
	if !found {
		mdPath := urlPath + ".markdown"
		if cached, ok := cache[mdPath]; ok {
			content = cached.Content
			title = filepath.Base(mdPath)
			currentPath = mdPath
			found = true
		}
	}

	// 4. Root path → try README.md / index.md
	if !found && urlPath == "" {
		for _, name := range []string{"README.md", "readme.md", "Readme.md", "index.md", "INDEX.md"} {
			if cached, ok := cache[name]; ok {
				content = cached.Content
				title = name
				currentPath = name
				found = true
				break
			}
		}
	}

	// 5. Directory path → try README.md / index.md inside it
	if !found && urlPath != "" {
		for _, name := range []string{"README.md", "readme.md", "Readme.md", "index.md", "INDEX.md"} {
			dirPath := urlPath + "/" + name
			if cached, ok := cache[dirPath]; ok {
				content = cached.Content
				title = name
				currentPath = dirPath
				found = true
				break
			}
		}
	}

	// 6. Show directory listing
	if !found {
		listing := buildDirectoryListing(urlPath, cache)
		if listing != "" {
			content = listing
			if urlPath == "" {
				title = filepath.Base(baseDir)
			} else {
				title = filepath.Base(urlPath)
			}
			found = true
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	// Rewrite internal .md links to clean URLs
	content = rewriteMdLinks(content)

	data := views.ServeData{
		Title:           title,
		Content:         content,
		ThemeName:       serveTheme,
		AutoReload:      serveReload,
		IsDirectoryMode: true,
		CurrentPath:     currentPath,
		Files:           tree,
	}
	templ.Handler(views.ServePage(data)).ServeHTTP(w, r)
}

// scanAndRenderDirectory walks baseDir and renders every .md / .markdown file.
func scanAndRenderDirectory() error {
	newCache := make(map[string]*CachedFile)
	var latestMod time.Time

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != baseDir {
				return filepath.SkipDir
			}
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		lower := strings.ToLower(info.Name())
		if !strings.HasSuffix(lower, ".md") && !strings.HasSuffix(lower, ".markdown") {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		content, err := renderer.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read %s: %v\n", relPath, err)
			return nil
		}

		rendered, err := renderer.Render(renderer.RenderOptions{
			Input:        content,
			Autolink:     true,
			Theme:        serveTheme,
			Width:        serveWidth,
			OutputFormat: "html",
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not render %s: %v\n", relPath, err)
			return nil
		}

		newCache[relPath] = &CachedFile{
			Content: rendered,
			ModTime: info.ModTime(),
		}
		if info.ModTime().After(latestMod) {
			latestMod = info.ModTime()
		}

		return nil
	})
	if err != nil {
		return err
	}

	tree := buildFileTree(newCache)

	fileCacheMu.Lock()
	fileCache = newCache
	fileTree = tree
	globalModTime = latestMod
	fileCacheMu.Unlock()

	return nil
}

// renderSingleCachedFile re-renders one file in the cache (for live-reload).
func renderSingleCachedFile(relPath string) error {
	absPath := filepath.Join(baseDir, filepath.FromSlash(relPath))

	stat, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	content, err := renderer.ReadFile(absPath)
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

	fileCacheMu.Lock()
	fileCache[relPath] = &CachedFile{
		Content: rendered,
		ModTime: stat.ModTime(),
	}
	if stat.ModTime().After(globalModTime) {
		globalModTime = stat.ModTime()
	}
	fileCacheMu.Unlock()

	return nil
}

// ====================================================================
// FILE TREE BUILDER
// ====================================================================

// dirNode is an intermediate structure used to build the file tree.
type dirNode struct {
	entry    views.FileEntry
	children map[string]*dirNode
	files    []views.FileEntry
}

func buildFileTree(cache map[string]*CachedFile) []views.FileEntry {
	paths := make([]string, 0, len(cache))
	for p := range cache {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	root := &dirNode{children: make(map[string]*dirNode)}

	for _, p := range paths {
		parts := strings.Split(p, "/")
		current := root
		for i, part := range parts {
			if i == len(parts)-1 {
				// Leaf file
				current.files = append(current.files, views.FileEntry{
					Name:  part,
					Path:  p,
					IsDir: false,
				})
			} else {
				// Directory – find or create
				child, ok := current.children[part]
				if !ok {
					child = &dirNode{
						entry: views.FileEntry{
							Name:  part,
							Path:  strings.Join(parts[:i+1], "/"),
							IsDir: true,
						},
						children: make(map[string]*dirNode),
					}
					current.children[part] = child
				}
				current = child
			}
		}
	}

	return convertTree(root)
}

func convertTree(node *dirNode) []views.FileEntry {
	var result []views.FileEntry

	// Directories first, sorted alphabetically
	dirNames := make([]string, 0, len(node.children))
	for name := range node.children {
		dirNames = append(dirNames, name)
	}
	sort.Strings(dirNames)

	for _, name := range dirNames {
		child := node.children[name]
		entry := child.entry
		entry.Children = convertTree(child)
		result = append(result, entry)
	}

	// Files second, sorted alphabetically
	sort.Slice(node.files, func(i, j int) bool {
		return node.files[i].Name < node.files[j].Name
	})
	result = append(result, node.files...)

	return result
}

// ====================================================================
// DIRECTORY LISTING (HTML)
// ====================================================================

func buildDirectoryListing(dirPath string, cache map[string]*CachedFile) string {
	type entry struct {
		name  string
		path  string
		isDir bool
	}

	seen := make(map[string]bool)
	var entries []entry

	prefix := dirPath
	if prefix != "" {
		prefix += "/"
	}

	for p := range cache {
		var rel string
		if prefix == "" {
			rel = p
		} else if strings.HasPrefix(p, prefix) {
			rel = p[len(prefix):]
		} else {
			continue
		}

		parts := strings.SplitN(rel, "/", 2)
		name := parts[0]

		if seen[name] {
			continue
		}
		seen[name] = true

		if len(parts) > 1 {
			entries = append(entries, entry{
				name:  name,
				path:  prefix + name,
				isDir: true,
			})
		} else {
			entries = append(entries, entry{
				name:  name,
				path:  prefix + name,
				isDir: false,
			})
		}
	}

	if len(entries) == 0 {
		return ""
	}

	// Sort: dirs first, then files, each alphabetical
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].isDir != entries[j].isDir {
			return entries[i].isDir
		}
		return entries[i].name < entries[j].name
	})

	var sb strings.Builder
	sb.WriteString(`<div class="space-y-6">`)

	// Breadcrumb navigation
	if dirPath != "" {
		sb.WriteString(`<nav class="flex items-center gap-1.5 text-sm text-muted-foreground mb-2">`)
		sb.WriteString(`<a href="/" class="hover:text-foreground transition-colors font-medium">Home</a>`)
		parts := strings.Split(dirPath, "/")
		for i, part := range parts {
			sb.WriteString(` <span class="text-muted-foreground/50">/</span> `)
			if i < len(parts)-1 {
				link := strings.Join(parts[:i+1], "/")
				sb.WriteString(fmt.Sprintf(`<a href="/%s" class="hover:text-foreground transition-colors">%s</a>`, link, part))
			} else {
				sb.WriteString(fmt.Sprintf(`<span class="text-foreground font-medium">%s</span>`, part))
			}
		}
		sb.WriteString(`</nav>`)
	}

	sb.WriteString(`<div class="grid gap-2">`)

	for _, e := range entries {
		icon := `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-muted-foreground shrink-0"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 9H8"/><path d="M16 13H8"/><path d="M16 17H8"/></svg>`
		href := "/" + e.path
		if e.isDir {
			icon = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-muted-foreground shrink-0"><path d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z"/></svg>`
			href += "/"
		}

		sb.WriteString(fmt.Sprintf(
			`<a href="%s" class="flex items-center gap-3 p-3 rounded-lg border bg-card hover:bg-accent hover:text-accent-foreground transition-colors group">%s<span class="font-medium">%s</span></a>`,
			href, icon, e.name))
	}

	sb.WriteString(`</div></div>`)
	return sb.String()
}

// ====================================================================
// LINK REWRITING
// ====================================================================

var mdLinkRegex = regexp.MustCompile(`href="([^"]*?)(\.md|\.markdown)(#[^"]*)?"`)

// rewriteMdLinks converts internal .md links to clean URLs for the server.
func rewriteMdLinks(content string) string {
	return mdLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Skip external URLs
		if strings.Contains(match, "://") {
			return match
		}
		match = strings.Replace(match, ".markdown#", "#", 1)
		match = strings.Replace(match, `.markdown"`, `"`, 1)
		match = strings.Replace(match, ".md#", "#", 1)
		match = strings.Replace(match, `.md"`, `"`, 1)
		return match
	})
}

// ====================================================================
// DIRECTORY WATCHER
// ====================================================================

func startDirectoryWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Watch all directories recursively
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") && path != baseDir {
				return filepath.SkipDir
			}
			if skipDirs[name] {
				return filepath.SkipDir
			}
			watcher.Add(path)
		}
		return nil
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Watch newly created directories
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, statErr := os.Stat(event.Name); statErr == nil && info.IsDir() {
					watcher.Add(event.Name)
					continue
				}
			}

			// Only process .md / .markdown files
			lower := strings.ToLower(event.Name)
			if !strings.HasSuffix(lower, ".md") && !strings.HasSuffix(lower, ".markdown") {
				continue
			}

			time.Sleep(200 * time.Millisecond) // Debounce

			relPath, relErr := filepath.Rel(baseDir, event.Name)
			if relErr != nil {
				continue
			}
			relPath = filepath.ToSlash(relPath)

			if event.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
				// File removed
				fileCacheMu.Lock()
				delete(fileCache, relPath)
				fileTree = buildFileTree(fileCache)
				globalModTime = time.Now()
				fileCacheMu.Unlock()

				if verbose {
					fmt.Printf("🗑️  File removed: %s at %s\n", relPath, time.Now().Format("15:04:05"))
				}
			} else if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				// File modified or created
				if renderErr := renderSingleCachedFile(relPath); renderErr != nil {
					fmt.Fprintf(os.Stderr, "Render error for %s: %v\n", relPath, renderErr)
				} else if verbose {
					fmt.Printf("📝 File updated: %s at %s\n", relPath, time.Now().Format("15:04:05"))
				}

				// Rebuild tree if a new file was created
				if event.Op&fsnotify.Create != 0 {
					fileCacheMu.Lock()
					fileTree = buildFileTree(fileCache)
					fileCacheMu.Unlock()
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

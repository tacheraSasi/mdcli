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

type PreviewData struct {
	Title      string
	Content    template.HTML
	Theme      themes.Theme
	AutoReload bool
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} – mdcli Preview</title>
    <!-- highlight.js for syntax highlighting -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
    <!-- Font Awesome for icons -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        /* Theme variables – injected from backend */
        :root {
            --bg: {{.Theme.Colors.Background}};
            --text: {{.Theme.Colors.Text}};
            --header: {{.Theme.Colors.Header}};
            --link: {{.Theme.Colors.Link}};
            --accent: {{.Theme.Colors.Accent}};
            --secondary: {{.Theme.Colors.Secondary}};
            --border: color-mix(in srgb, var(--text) 15%, transparent);
            --code-bg: color-mix(in srgb, var(--text) 8%, transparent);
            --toc-bg: color-mix(in srgb, var(--bg) 95%, black);
            --shadow: 0 4px 12px rgba(0,0,0,0.1);
        }

        /* Dark mode overrides (toggle via class) */
        body.dark {
            --bg: #1e1e1e;
            --text: #e0e0e0;
            --header: #bb86fc;
            --link: #8ab4f8;
            --accent: #03dac6;
            --secondary: #03dac6;
            --border: #333;
            --code-bg: #2d2d2d;
            --toc-bg: #252525;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.7;
            transition: background 0.2s ease, color 0.2s ease;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }

        /* Main layout */
        .container {
            display: flex;
            max-width: 1400px;
            margin: 0 auto;
            padding: 1rem;
            gap: 2rem;
            flex: 1;
        }

        /* Sidebar (TOC) */
        .toc {
            width: 260px;
            flex-shrink: 0;
            position: sticky;
            top: 1rem;
            align-self: start;
            background: var(--toc-bg);
            border-radius: 12px;
            padding: 1.2rem 1rem;
            max-height: calc(100vh - 2rem);
            overflow-y: auto;
            border: 1px solid var(--border);
            box-shadow: var(--shadow);
            font-size: 0.9rem;
        }

        .toc h3 {
            margin-bottom: 1rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.5rem;
            color: var(--header);
        }

        .toc ul {
            list-style: none;
        }

        .toc li {
            margin: 0.5rem 0;
            padding-left: 0.5rem;
            border-left: 2px solid transparent;
        }

        .toc li.active {
            border-left-color: var(--accent);
            font-weight: 500;
        }

        .toc a {
            color: var(--text);
            text-decoration: none;
            display: block;
            transition: color 0.1s;
            word-break: break-word;
        }

        .toc a:hover {
            color: var(--link);
        }

        .toc .h2 { margin-left: 1rem; }
        .toc .h3 { margin-left: 2rem; }
        .toc .h4 { margin-left: 3rem; }

        /* Main content */
        .content {
            flex: 1;
            max-width: 800px;
            margin: 0 auto;
            width: 100%;
        }

        /* Header area */
        .doc-header {
            margin-bottom: 2.5rem;
            padding-bottom: 1rem;
            border-bottom: 2px solid var(--border);
        }

        .doc-header h1 {
            font-size: 2.5rem;
            color: var(--header);
            margin-bottom: 0.25rem;
        }

        .doc-meta {
            display: flex;
            gap: 1.5rem;
            color: var(--secondary);
            font-size: 0.9rem;
            align-items: center;
            flex-wrap: wrap;
        }

        /* Floating action bar */
        .action-bar {
            position: fixed;
            bottom: 2rem;
            right: 2rem;
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
            z-index: 100;
        }

        .action-btn {
            width: 48px;
            height: 48px;
            border-radius: 50%;
            background: var(--accent);
            color: white;
            border: none;
            cursor: pointer;
            box-shadow: var(--shadow);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.2rem;
            transition: transform 0.2s, background 0.2s;
            backdrop-filter: blur(4px);
        }

        .action-btn:hover {
            transform: scale(1.1);
            background: var(--link);
        }

        /* Live indicator */
        .live-indicator {
            position: fixed;
            top: 1rem;
            right: 1rem;
            background: var(--accent);
            color: white;
            padding: 0.5rem 1.2rem;
            border-radius: 40px;
            font-size: 0.85rem;
            font-weight: 500;
            box-shadow: var(--shadow);
            backdrop-filter: blur(4px);
            z-index: 200;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .live-indicator i {
            font-size: 0.7rem;
            animation: pulse 1.5s infinite;
        }

        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.4; }
            100% { opacity: 1; }
        }

        /* Typography */
        h1, h2, h3, h4, h5, h6 {
            color: var(--header);
            margin-top: 2rem;
            margin-bottom: 1rem;
            font-weight: 600;
            line-height: 1.3;
        }

        h1 { font-size: 2.2rem; margin-top: 0; }
        h2 { font-size: 1.8rem; border-bottom: 1px solid var(--border); padding-bottom: 0.3rem; }
        h3 { font-size: 1.5rem; }
        h4 { font-size: 1.3rem; }

        p, ul, ol {
            margin-bottom: 1.2rem;
        }

        a {
            color: var(--link);
            text-decoration: none;
            border-bottom: 1px dotted currentColor;
        }

        a:hover {
            border-bottom-style: solid;
        }

        /* Code blocks */
        pre {
            background: var(--code-bg);
            border-radius: 8px;
            padding: 1.2rem;
            overflow-x: auto;
            position: relative;
            border: 1px solid var(--border);
            margin: 1.5rem 0;
        }

        code {
            font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
            font-size: 0.9em;
        }

        p code, li code {
            background: var(--code-bg);
            padding: 0.2em 0.4em;
            border-radius: 4px;
            border: 1px solid var(--border);
        }

        /* Copy button */
        .copy-btn {
            position: absolute;
            top: 0.5rem;
            right: 0.5rem;
            background: var(--bg);
            border: 1px solid var(--border);
            color: var(--text);
            border-radius: 4px;
            padding: 0.2rem 0.5rem;
            font-size: 0.8rem;
            cursor: pointer;
            opacity: 0;
            transition: opacity 0.2s;
        }

        pre:hover .copy-btn {
            opacity: 1;
        }

        /* Blockquotes */
        blockquote {
            margin: 1.5rem 0;
            padding: 0.5rem 1.5rem;
            border-left: 4px solid var(--accent);
            background: var(--code-bg);
            border-radius: 0 8px 8px 0;
            font-style: italic;
        }

        /* Images */
        img {
            max-width: 100%;
            border-radius: 8px;
            box-shadow: var(--shadow);
        }

        /* Responsive */
        @media (max-width: 900px) {
            .container {
                flex-direction: column;
            }
            .toc {
                width: 100%;
                position: relative;
                max-height: none;
            }
            .action-bar {
                bottom: 1rem;
                right: 1rem;
            }
        }

        /* Toggle for TOC on mobile */
        .toc-collapse {
            display: none;
        }
        @media (max-width: 600px) {
            .toc.collapsed {
                display: none;
            }
            .toc-collapse {
                display: block;
                margin-bottom: 1rem;
            }
        }
    </style>
    {{if .AutoReload}}
    <script>
        let lastModified = 0;
        function checkForUpdates() {
            fetch('/status')
                .then(r => r.json())
                .then(data => {
                    if (data.lastModified !== lastModified && lastModified !== 0) {
                        location.reload();
                    }
                    lastModified = data.lastModified;
                })
                .catch(() => {});
        }
        setInterval(checkForUpdates, 1000);
        checkForUpdates();
    </script>
    {{end}}
</head>
<body>
    <!-- Live indicator -->
    <div class="live-indicator">
        <i class="fas fa-circle"></i> LIVE
        <span id="last-updated"></span>
    </div>

    <!-- Floating action buttons -->
    <div class="action-bar">
        <button class="action-btn" id="theme-toggle" title="Toggle dark/light mode">
            <i class="fas fa-moon"></i>
        </button>
        <button class="action-btn" id="toc-toggle" title="Toggle table of contents">
            <i class="fas fa-list"></i>
        </button>
        <button class="action-btn" id="back-to-top" title="Back to top">
            <i class="fas fa-arrow-up"></i>
        </button>
    </div>

    <div class="container">
        <!-- Table of Contents (generated by JS) -->
        <nav class="toc" id="toc">
            <h3><i class="fas fa-list-ul"></i> Contents</h3>
            <ul id="toc-list"></ul>
        </nav>

        <!-- Main content -->
        <main class="content">
            <header class="doc-header">
                <h1>{{.Title}}</h1>
                <div class="doc-meta">
                    <span><i class="far fa-file-alt"></i> mdcli v2</span>
                    <span><i class="far fa-clock"></i> <span id="file-mod-time">loading...</span></span>
                    <span><i class="fas fa-palette"></i> Theme: {{.Theme.Name}}</span>
                </div>
            </header>
            <article>
                {{.Content}}
            </article>
        </main>
    </div>

    <script>
        (function() {
            // ========== THEME TOGGLE ==========
            const themeToggle = document.getElementById('theme-toggle');
            const icon = themeToggle.querySelector('i');
            // Check local storage or system preference
            const storedTheme = localStorage.getItem('theme');
            if (storedTheme === 'dark') {
                document.body.classList.add('dark');
                icon.classList.remove('fa-moon');
                icon.classList.add('fa-sun');
            } else if (storedTheme === 'light') {
                document.body.classList.remove('dark');
                icon.classList.remove('fa-sun', 'fa-moon');
                icon.classList.add('fa-moon');
            } else {
                // Auto: check prefers-color-scheme
                if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                    document.body.classList.add('dark');
                    icon.classList.remove('fa-moon');
                    icon.classList.add('fa-sun');
                }
            }

            themeToggle.addEventListener('click', () => {
                document.body.classList.toggle('dark');
                const isDark = document.body.classList.contains('dark');
                icon.classList.toggle('fa-moon', !isDark);
                icon.classList.toggle('fa-sun', isDark);
                localStorage.setItem('theme', isDark ? 'dark' : 'light');
            });

            // ========== TABLE OF CONTENTS ==========
            const tocList = document.getElementById('toc-list');
            const headings = document.querySelectorAll('article h1, article h2, article h3, article h4, article h5, article h6');
            const tocLinks = [];

            if (headings.length > 0) {
                headings.forEach((heading, index) => {
                    // Ensure each heading has an id for linking
                    if (!heading.id) {
                        heading.id = heading.tagName.toLowerCase() + '-' + index;
                    }
                    const li = document.createElement('li');
                    li.className = heading.tagName.toLowerCase(); // h1, h2, etc.
                    const a = document.createElement('a');
                    a.href = '#' + heading.id;
                    a.textContent = heading.textContent;
                    li.appendChild(a);
                    tocList.appendChild(li);
                    tocLinks.push({ link: a, heading: heading });
                });
            } else {
                tocList.innerHTML = '<li><em>No headings</em></li>';
            }

            // Highlight active TOC item on scroll
            function setActiveTOC() {
                const scrollPos = window.scrollY + 80;
                let current = null;
                for (let i = headings.length - 1; i >= 0; i--) {
                    const heading = headings[i];
                    if (heading.offsetTop <= scrollPos) {
                        current = heading;
                        break;
                    }
                }
                document.querySelectorAll('.toc li').forEach(li => li.classList.remove('active'));
                if (current) {
                    const activeLink = document.querySelector('.toc a[href="#${current.id}"]');
                    if (activeLink) activeLink.parentElement.classList.add('active');
                }
            }
            window.addEventListener('scroll', setActiveTOC);
            setActiveTOC();

            // ========== COPY CODE BUTTONS ==========
            document.querySelectorAll('pre').forEach(pre => {
                const btn = document.createElement('button');
                btn.className = 'copy-btn';
                btn.innerHTML = '<i class="far fa-copy"></i> Copy';
                btn.addEventListener('click', () => {
                    const code = pre.querySelector('code');
                    const text = code.innerText;
                    navigator.clipboard.writeText(text).then(() => {
                        btn.innerHTML = '<i class="fas fa-check"></i> Copied!';
                        setTimeout(() => {
                            btn.innerHTML = '<i class="far fa-copy"></i> Copy';
                        }, 2000);
                    });
                });
                pre.style.position = 'relative';
                pre.appendChild(btn);
            });

            // ========== BACK TO TOP ==========
            document.getElementById('back-to-top').addEventListener('click', () => {
                window.scrollTo({ top: 0, behavior: 'smooth' });
            });

            // ========== TOGGLE TOC ON MOBILE ==========
            const toc = document.getElementById('toc');
            const tocToggle = document.getElementById('toc-toggle');
            tocToggle.addEventListener('click', () => {
                toc.classList.toggle('collapsed');
            });

            // ========== UPDATE LAST MODIFIED TIME ==========
            function updateModTime() {
                fetch('/status')
                    .then(r => r.json())
                    .then(data => {
                        const d = new Date(data.lastModified * 1000);
                        document.getElementById('file-mod-time').textContent = d.toLocaleTimeString();
                    })
                    .catch(() => {});
            }
            updateModTime();
            setInterval(updateModTime, 5000);

            // ========== SYNTAX HIGHLIGHTING ==========
            // Load highlight.js and apply to all code blocks
            var script = document.createElement('script');
            script.src = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js';
            script.onload = function() {
                hljs.highlightAll();
            };
            document.head.appendChild(script);
        })();
    </script>
</body>
</html>`

var (
	currentFile     string
	lastModTime     time.Time
	cachedContent   string
	previewTemplate *template.Template
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

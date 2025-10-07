# 🚀 MDCLI v2.0 - Advanced Markdown CLI Processor

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/badge/Release-v2.0.0-brightgreen.svg)](https://github.com/tacheraSasi/mdcli/releases)

**mdcli** is a powerful, feature-rich command-line tool for processing and rendering Markdown files with advanced capabilities including live preview, batch processing, multiple output formats, and extensive customization options.

## ✨ Features

### 🎨 **Multiple Output Formats**

- **Terminal**: Beautiful terminal rendering with syntax highlighting
- **HTML**: Clean HTML output with custom themes
- **PDF**: Export to PDF format (planned)
- **Plain Text**: Strip HTML for plain text output

### 🌈 **Theme Support**

- **5 Built-in Themes**: Dracula, GitHub, Monokai, Solarized, Nord
- **Customizable Colors**: Override theme colors via configuration
- **Syntax Highlighting**: Advanced code highlighting with multiple styles

### 🔄 **Live Features**

- **Watch Mode**: Auto-regenerate on file changes
- **Live Preview Server**: Browser-based live preview with auto-reload
- **Interactive Mode**: Real-time Markdown editor in terminal

### ⚡ **Performance & Productivity**

- **Batch Processing**: Process entire directories concurrently
- **Progress Bars**: Visual feedback for long operations
- **Concurrent Workers**: Configurable parallel processing
- **Caching**: Intelligent caching for better performance

### 🛠️ **Advanced Markdown**

- **Math Equations**: LaTeX math rendering with MathJax
- **Mermaid Diagrams**: Support for flowcharts and diagrams
- **GitHub Flavored Markdown**: Tables, task lists, strikethrough
- **Auto-linking**: Automatic URL and email detection

### ⚙️ **Configuration**

- **YAML Config Files**: Persistent settings in `~/.mdcli.yaml`
- **Environment Variables**: Override settings with `MDCLI_*` vars
- **Command-line Flags**: Override any setting on-the-fly

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/tacheraSasi/mdcli.git
cd mdcli

# Build for your platform
make build_linux    # Linux
make build_mac      # macOS
make build_windows  # Windows

# Or build directly with Go
go build -o mdcli
```

### Basic Usage

```bash
# Render a file to terminal
./mdcli render README.md

# Render with custom theme
./mdcli render --theme=github README.md

# Output to HTML
./mdcli render --format=html --output=output.html README.md

# Start live preview server
./mdcli serve README.md --port=8080

# Watch for changes
./mdcli watch README.md --output=preview.html

# Interactive mode
./mdcli interactive

# Batch process directory
./mdcli batch ./docs --output=./dist --recursive

# View available themes
./mdcli themes
```

## 📋 Command Reference

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `render` | Render Markdown files | `mdcli render file.md` |
| `serve` | Start live preview server | `mdcli serve file.md` |
| `watch` | Watch files for changes | `mdcli watch file.md` |
| `batch` | Process multiple files | `mdcli batch ./docs` |
| `interactive` | Interactive editor mode | `mdcli interactive` |
| `themes` | List available themes | `mdcli themes` |
| `config` | Manage configuration | `mdcli config show` |

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `~/.mdcli.yaml` |
| `--verbose` | Verbose output | `false` |
| `--theme` | Syntax theme | `dracula` |
| `--format` | Output format | `terminal` |
| `--width` | Terminal width | `80` |

### Render Command Options

```bash
mdcli render [files...] [flags]

Flags:
  -f, --format string   Output format (terminal, html, pdf, text)
  -o, --output string   Output file path
  -t, --theme string    Syntax highlighting theme
  -w, --width int       Terminal width for formatting
      --autolink        Enable automatic link detection (default true)
      --progress        Show progress bar for multiple files
```

### Serve Command Options

```bash
mdcli serve [file] [flags]

Flags:
  -p, --port int        Port to serve on (default 8080)
  -b, --bind string     Bind address (default "localhost")
      --auto-reload     Enable auto-reload on file changes (default true)
  -t, --theme string    Theme for HTML output (default "github")
```

### Batch Command Options

```bash
mdcli batch [directory] [flags]

Flags:
  -o, --output string      Output directory
  -f, --format string      Output format (default "html")
  -r, --recursive          Process subdirectories recursively
  -c, --concurrent int     Number of concurrent workers (default 4)
  -e, --ext string         Output file extension (default ".html")
```

## ⚙️ Configuration

### Configuration File

Create a configuration file at `~/.mdcli.yaml`:

```bash
# Initialize default config
mdcli config init

# View current configuration
mdcli config show
```

### Sample Configuration

```yaml
# Default theme for syntax highlighting
theme: dracula

# Default output format
output_format: terminal

# Default terminal width
width: 80

# Enable automatic link detection
autolink: true

# Batch processing settings
batch:
  concurrent_workers: 4
  output_extension: .html
  recursive: true

# Server settings
serve:
  port: 8080
  bind: localhost
  auto_reload: true

# Rendering preferences
render:
  show_progress: true
  include_metadata: false
```

### Environment Variables

Override any setting with environment variables:

```bash
export MDCLI_THEME=github
export MDCLI_WIDTH=100
export MDCLI_FORMAT=html

mdcli render README.md  # Uses above settings
```

## 🎨 Themes

mdcli includes 5 beautiful built-in themes:

| Theme | Description | Best For |
|-------|-------------|----------|
| **Dracula** | Dark theme with purple accents | Dark mode lovers |
| **GitHub** | Light theme inspired by GitHub | Documentation |
| **Monokai** | Dark theme with vibrant colors | Code highlighting |
| **Solarized** | Balanced light/dark theme | Eye comfort |
| **Nord** | Arctic, north-bluish palette | Modern aesthetic |

```bash
# List all available themes
mdcli themes

# Use a specific theme
mdcli render --theme=nord README.md

# View theme details
mdcli themes --verbose
```

## 🔧 Advanced Features

### Math Equations

mdcli supports LaTeX math equations via MathJax:

```markdown
Inline math: $E = mc^2$

Block math:
$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$
```

### Mermaid Diagrams

Create flowcharts and diagrams:

```markdown
\```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
\```
```

### Live Preview

Start a live preview server for real-time editing:

```bash
# Start server
mdcli serve README.md --port=8080

# Open browser to http://localhost:8080
# Edit README.md and see changes instantly
```

### Interactive Mode

Use the interactive REPL for quick Markdown testing:

```bash
mdcli interactive

# Enter Markdown content, press Enter twice to render
# Type 'exit' to quit
```

## 🏗️ Building from Source

### Prerequisites

- Go 1.24 or later
- Make (optional, for using Makefile)

### Build Commands

```bash
# Build for current platform
go build -o mdcli

# Build for all platforms
make build

# Build optimized binaries
make build_linux    # Linux AMD64
make build_mac      # macOS AMD64  
make build_windows  # Windows AMD64
make build_android  # Android ARM64

# Clean build artifacts
make clean

# Install dependencies
make dependencies
```

### Cross-compilation

```bash
# Manual cross-compilation examples
GOOS=linux GOARCH=amd64 go build -o mdcli-linux
GOOS=darwin GOARCH=amd64 go build -o mdcli-darwin
GOOS=windows GOARCH=amd64 go build -o mdcli.exe
```

## 📈 Performance Tips

1. **Use batch mode** for processing multiple files
2. **Enable progress bars** for visual feedback: `--progress`
3. **Adjust concurrent workers**: `--concurrent=8` for more CPU cores
4. **Use appropriate themes** for your terminal
5. **Set optimal width**: `--width=120` for wide terminals

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone repository
git clone https://github.com/tacheraSasi/mdcli.git
cd mdcli

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build and test
go build && ./mdcli render README.md
```

## 📝 Changelog

### v2.0.0 (Latest)

- ✨ Complete rewrite with Cobra CLI framework
- 🎨 Multiple output formats (HTML, PDF, plain text)
- 🌈 5 built-in themes with customization
- 🔄 Live preview server with auto-reload
- ⚡ Batch processing with concurrent workers
- 🛠️ Advanced Markdown features (Math, Mermaid)
- ⚙️ Comprehensive configuration system
- 📱 Interactive mode for real-time editing

### v1.0.0

- 📄 Basic Markdown rendering
- 🖥️ Terminal output only
- 📁 Single file processing

## 🐛 Troubleshooting

### Common Issues

**Issue**: Command not found
```bash
# Solution: Ensure binary is in PATH or use full path
./mdcli render file.md
```

**Issue**: Theme not working
```bash
# Solution: Check available themes
mdcli themes
mdcli render --theme=github file.md
```

**Issue**: Server won't start
```bash
# Solution: Use different port
mdcli serve file.md --port=3000
```

**Issue**: Config file not found
```bash
# Solution: Initialize config
mdcli config init
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [go-term-markdown](https://github.com/MichaelMure/go-term-markdown) - Terminal rendering

## 📞 Support

- 📧 Email: [tacherasasi@gmail.com](mailto:tacherasasi@gmail.com)
- 🐛 Issues: [GitHub Issues](https://github.com/tacheraSasi/mdcli/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/tacheraSasi/mdcli/discussions)

---

<p align="center">
  <b>Made with ❤️ by <a href="https://github.com/tacheraSasi">Tachera Sasi</a></b>
</p>

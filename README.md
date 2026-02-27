# MDCLI v2.0 - Advanced Markdown CLI Processor

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/badge/Release-v2.0.0-brightgreen.svg)](https://github.com/tacheraSasi/mdcli/releases)

**mdcli** is a powerful, feature-rich command-line tool for processing and rendering Markdown files. It provides advanced capabilities including live preview, batch processing, multiple output formats, and extensive customization options.

> **Backward Compatibility**: mdcli v2.0 maintains full backward compatibility. You can still use `mdcli file.md` directly without the `render` subcommand.

## Features

### Multiple Output Formats

- **Terminal**: High-quality terminal rendering with syntax highlighting
- **HTML**: Clean HTML output with customizable themes
- **PDF**: Export to PDF format (planned)
- **Plain Text**: Strip formatting for plain text output

### Theme Support

- **Built-in Themes**: Dracula, GitHub, Monokai, Solarized, and Nord
- **Customizable Colors**: Override theme colors via configuration
- **Syntax Highlighting**: Advanced code highlighting with multiple styles

### Live Features

- **Watch Mode**: Automatic regeneration on file changes
- **Live Preview Server**: Browser-based live preview with auto-reload
- **Interactive Mode**: Real-time Markdown editor within the terminal

### Performance and Productivity

- **Batch Processing**: Process entire directories concurrently
- **Progress Monitoring**: Visual feedback for long-running operations
- **Concurrent Workers**: Configurable parallel processing
- **Intelligent Caching**: Optimized performance for repeated tasks

### Advanced Markdown Support

- **Math Equations**: LaTeX math rendering with MathJax
- **Mermaid Diagrams**: Support for flowcharts and diagrams
- **GitHub Flavored Markdown**: Tables, task lists, and strikethrough
- **Auto-linking**: Automatic URL and email detection

### Configuration

- **YAML Configuration**: Persistent settings via `~/.mdcli.yaml`
- **Environment Variables**: Override settings with `MDCLI_*` variables
- **Command-line Flags**: On-the-fly overrides for any setting

## Quick Start

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
# Simple rendering (backward compatible)
./mdcli README.md

# With flags (backward compatible)
./mdcli README.md --theme=github --format=html

# Explicit render command
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

## Command Reference

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

## Configuration

### Configuration File

Manage your settings using the built-in configuration tool:

```bash
# Initialize default config
mdcli config init

# View current configuration
mdcli config show
```

### Sample Configuration (`~/.mdcli.yaml`)

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

Settings can be overridden using environment variables:

```bash
export MDCLI_THEME=github
export MDCLI_WIDTH=100
export MDCLI_FORMAT=html

mdcli render README.md
```

## Themes

mdcli includes several professionally designed built-in themes:

| Theme | Description | Best For |
|-------|-------------|----------|
| **Dracula** | Dark theme with purple accents | Dark mode environments |
| **GitHub** | Light theme inspired by GitHub | Documentation |
| **Monokai** | Dark theme with vibrant colors | Code-heavy documents |
| **Solarized** | Balanced light/dark theme | Visual comfort |
| **Nord** | Arctic, north-bluish palette | Modern interfaces |

```bash
# List all available themes
mdcli themes

# Use a specific theme
mdcli render --theme=nord README.md
```

## Advanced Features

### Math Equations

Supports LaTeX math equations via MathJax:

```markdown
Inline math: $E = mc^2$

Block math:
$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$
```

### Mermaid Diagrams

Native support for flowcharts and diagrams:

```markdown
\```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
\```
```

### Live Preview

Real-time browser-based preview:

```bash
# Start server
mdcli serve README.md --port=8080

# The server will automatically reload when files are saved.
```

## Building from Source

### Prerequisites

- Go 1.24 or later
- Make (optional)

### Build Commands

```bash
# Build for the current platform
go build -o mdcli

# Build for all platforms via Makefile
make build

# Optimized binaries
make build_linux
make build_mac
make build_windows
```

## Contributing

Contributions are welcome. Please refer to the [Contributing Guidelines](CONTRIBUTING.md) for further details.

## Changelog

### v2.0.0
- Complete migration to the Cobra CLI framework
- Support for multiple output formats (HTML, PDF, Text)
- Enhanced theme system with 5 built-in options
- Integrated live preview server with auto-reload
- Concurrent batch processing
- Support for MathJax and Mermaid diagrams
- Comprehensive YAML-based configuration

### v1.0.0
- Initial release with basic Markdown rendering
- Terminal output support

## Support

- **Email**: [tacherasasi@gmail.com](mailto:tacherasasi@gmail.com)
- **Issues**: [GitHub Issues](https://github.com/tacheraSasi/mdcli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tacheraSasi/mdcli/discussions)

---

<p align="center">
  Developed by <a href="https://github.com/tacheraSasi">Tachera Sasi</a>
</p>

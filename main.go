package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tacheraSasi/mdcli/renderer"
)

var (
	version = "1.0.0"
	style   = flag.String("style", "dark", "Style to use for rendering")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var input string
	var err error

	if flag.NArg() == 0 {
		// Read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		input = string(content)
	} else {
		filename := flag.Arg(0)
		if !strings.HasSuffix(filename, ".md") {
			return fmt.Errorf("file must have a .md extension")
		}
		content, err := renderer.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		input = content
	}

	rendered, err := renderer.Render(input, *style)
	if err != nil {
		return fmt.Errorf("error rendering markdown: %w", err)
	}

	fmt.Println(rendered)
	return nil
}

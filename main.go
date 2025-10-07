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
	version  = "1.0.0"
	output   = flag.String("output", "", "Write rendered output to file (default: stdout)")
	autolink = flag.Bool("autolink", true, "Enable or disable autolinking")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [options] [file ...]\n", "mdcli")
		fmt.Fprintln(os.Stdout, "Options:")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "Examples:")
		fmt.Fprintln(os.Stdout, "  mdcli README.md")
		fmt.Fprintln(os.Stdout, "  mdcli --output=out.txt file1.md file2.md")
		fmt.Fprintln(os.Stdout, "  cat README.md | mdcli")
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
		rendered, err := renderer.Render(renderer.RenderOptions{
			Input:    input,
			Autolink: *autolink,
		})
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
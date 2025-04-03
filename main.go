package main

import (
	"flag"
	"fmt"

	"github.com/tacheraSasi/mdcli/renderer"
)

var VERSION string = "1"

func main() {
	filename := flag.String("file", "", "Path to the file")

	flag.Parse()

	if *filename == "" {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("No file provided")
			return
		} 
		*filename = args[0]
	}

	mdFile := *filename
	mdFileContent, err := renderer.ReadFile(mdFile)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	rendered, err := renderer.Render(mdFileContent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rendered)

}

package main

import (
	"flag"
	"fmt"
)

var VERSION string= "1"

func main() {
	filename := flag.String("file", "", "Path to the file")

	flag.Parse()

	if *filename == "" {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("No file provided")
			return
		}
	}
	
}
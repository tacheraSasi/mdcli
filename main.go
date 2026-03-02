package main

import "github.com/tacheraSasi/mdcli/cmd"

func main() {
	// Wire embedded assets into the serve command
	cmd.AssetsFS = EmbeddedAssets
	cmd.Execute()
}

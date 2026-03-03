package main

import "embed"

// EmbeddedAssets contains the CSS and JS files baked into the binary.
//
//go:embed assets/css/output.css assets/css/input.css assets/js/*.js assets/icons
var EmbeddedAssets embed.FS

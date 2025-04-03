package main // Changed package to main

import (
	"embed"
	_ "embed" // Import for side-effect of //go:embed
)

//go:embed all:templates
var EmbeddedTemplateFS embed.FS

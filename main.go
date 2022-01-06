package main

import (
	_ "embed"
	"mcmods/cmd"
)

//go:embed VERSION.txt
var version string

func main() {
	cmd.ToolVersion = version
	cmd.Execute()
}

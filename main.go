package main

import (
	_ "embed"
	"mcmods/cmd"
)

//go:embed VERSION.txt
var version string

func main() {
	cmd.Version = version
	cmd.Execute()
}

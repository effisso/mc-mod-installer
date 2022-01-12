package cmd

import (
	"github.com/spf13/cobra"
)

// DocURL is the homepage of the online documentation for this tool
var DocURL = "https://github.com/effisso/mc-mod-installer/tree/main/docs"

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Open the documentation homepage for this tool.",
	Long: `
No arguments. Simply opens a browser to:
https://github.com/effisso/mc-mod-installer/tree/main/docs`,
	Run: func(cmd *cobra.Command, args []string) {
		BrowserLauncher.Open(DocURL)
	},
}

func init() {
	RootCmd.AddCommand(docsCmd)
}

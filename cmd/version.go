package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// ToolVersion represents the version of this tool
	ToolVersion string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the running version of this tool.",
	Long: `
Print the running version of this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		printToUser(ToolVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

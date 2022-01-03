package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [resource]",
	Short: "List mods and server groups",
	Long: `
See
 $ list mods --help
or
 $ list groups --help
for more information on usage.`,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

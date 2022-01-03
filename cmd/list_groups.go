package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

// listGroupsCmd represents the group command
var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List server groups by name",
	Long: `
Prints out all the valid group names.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, group := range getServerModGroupNames(mc.ServerGroups) {
			printToUser(group)
		}
	},
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}

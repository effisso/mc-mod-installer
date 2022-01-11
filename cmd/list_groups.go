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
		groupNames := getServerModGroupNames(mc.ServerGroups)
		max := len(groupNames) - 1
		for i, group := range groupNames {
			if i == max {
				printToUser(group)
			} else {
				printLineToUser(group)
			}
		}
	},
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}

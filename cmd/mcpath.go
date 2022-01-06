package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	path *string
)

// pathCmd represents the loc command
var pathCmd = &cobra.Command{
	Use:   "mcpath",
	Short: "Get and set the path to the Minecraft install folder",
	Long: `
Mods must be installed in a specific place inside the Minecraft installation
directory. This command helps ensure the tool is using the right path for this
machine.

To print the location, call this command with no args. To set the path, use the
--set option:
 $ mcpath --set /absolute/path/to/.minecraft`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if *path == "" {
			printToUser(mc.GetInstallPath())
		} else {
			ViperInstance.Set(mc.InstallPathKey, path)
			err = ViperInstance.WriteConfig()
			if err == nil {
				printToUser("Path updated.")
			}
		}
		return
	},
}

// ResetPathVars is used for testing
func ResetPathVars() {
	*path = ""
}

func init() {
	ResetAddVars()

	RootCmd.AddCommand(pathCmd)

	path = pathCmd.Flags().String("set", "", "Used to set the path where Minecraft is installed: --set /path/to/folder")
}
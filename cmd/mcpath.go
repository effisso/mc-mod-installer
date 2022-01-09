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
directory. This command helps ensure the installer is using the right path for
this machine's Minecraft installation.

To print the location which the tool is attempting to use, call this command
with no args. To change the path this tool is using, include the --set option
and an absolute path:
 $ mcpath --set /absolute/path/to/.minecraft

Note: No validation is done on the provided path. Make sure it's correct!
Surround the path with double-quotes if it contains spaces. `,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if *path == "" {
			printToUser(mc.GetInstallPath())
		} else {
			ViperInstance.Set(mc.InstallPathKey, path)
			cobra.CheckErr(ViperInstance.WriteConfig())
			printToUser("Path updated.")
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

	path = pathCmd.Flags().String("set", "", "Used to set the absolute path where Minecraft is installed: --set /path/to/.minecraft")
}

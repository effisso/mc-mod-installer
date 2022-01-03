package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	path = newStrPtr()
)

// pathCmd represents the loc command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Get and set the path to the Minecraft install folder",
	Long: `
Mods must be installed in a specific place inside the Minecraft installation
directory. This command helps ensure the tool is using the right path for this
machine.

To print the location, call this command with no args. To set the path, use the
--set option:
 $ path --set /absolute/path/to/.minecraft`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if *path == "" {
			printToUser(ViperInstance.GetString(mc.InstallPathKey))
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

func init() {
	RootCmd.AddCommand(pathCmd)

	pathCmd.Flags().StringVar(path, "set", "", "Used to set the path where Minecraft is installed: --set /path/to/folder")
}

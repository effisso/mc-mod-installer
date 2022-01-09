package cmd

import (
	"mcmods/browser"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	BrowserLauncher = browser.NewLauncher()
)

// visitCmd represents the visit command
var visitCmd = &cobra.Command{
	Use:   "visit [mod name]",
	Short: "Open a browser to the mod's details URL",
	Long: `
Open a browser to the mod's details URL`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modName := args[0]
		cliMods := NameMapper.MapAllMods(UserModConfig.ClientMods)

		m := cliMods[modName]

		if m == nil {
			return mc.NewUnknownModError(modName)
		}

		BrowserLauncher.Open(m.DetailsURL)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(visitCmd)
}

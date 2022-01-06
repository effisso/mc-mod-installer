package cmd

import (
	"mcmods/browser"

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
	Run: func(cmd *cobra.Command, args []string) {
		modName := args[0]
		cliMods := NameMapper.MapAllMods(UserModConfig.ClientMods)

		err := NameValidator.ValidateModCliNames([]string{modName}, cliMods)
		cobra.CheckErr(err)

		m := cliMods[modName]

		BrowserLauncher.Open(m.DetailsURL)
	},
}

func init() {
	RootCmd.AddCommand(visitCmd)
}

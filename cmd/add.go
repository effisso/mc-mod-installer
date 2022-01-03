package cmd

import (
	"mcmods/input"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

const (
	AddFriendlyPromptText  = "What's the mod's name?\n> "
	AddCliNamePromptText   = "Globally unique name (lowercase letters and hyphens only; short yet descriptive)\n> "
	AddDescPromptText      = "Description of the mod (optional)\n> "
	AddDetailUrlPromptText = "Mod homepage/wiki URL\n> "
	AddDownloadPromptText  = "Desired package download URL\n> "
	AddGroupNamePromptText = "Server group\n> "
)

var (
	FriendlyPrompt    input.Prompt
	CliNamePrompt     input.Prompt
	DescPrompt        input.Prompt
	DetailsUrlPrompt  input.Prompt
	DownloadUrlPrompt input.Prompt
	GroupPrompt       input.Prompt

	ServerCfgSaver mc.ServerConfigSaver

	serverMod *bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [options]",
	Short: "Add a new mod to the installer",
	Long: `
Add a mod to the configuration to be installed and managed by this tool.
Typical usage is to add a client-side mod. The command can also add a
server-side mod when building new versions of the tool itself by using the
--server option.

All inputs for the mod information are collected interactively during
execution.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var friendlyName, cliName, desc, detUrl, dlUrl, groupName string

		out := cmd.OutOrStdout()
		in := cmd.InOrStdin()

		if friendlyName, err = FriendlyPrompt.GetInput(out, in); err != nil {
			return
		}

		if cliName, err = CliNamePrompt.GetInput(out, in); err != nil {
			return
		}

		if desc, err = DescPrompt.GetInput(out, in); err != nil {
			return
		}

		if detUrl, err = DetailsUrlPrompt.GetInput(out, in); err != nil {
			return
		}

		if dlUrl, err = DownloadUrlPrompt.GetInput(out, in); err != nil {
			return
		}

		mod := &mc.Mod{
			FriendlyName: friendlyName,
			CliName:      cliName,
			Description:  desc,
			DetailsUrl:   detUrl,
			LatestUrl:    dlUrl,
		}

		if *serverMod {
			if groupName, err = GroupPrompt.GetInput(out, in); err != nil {
				return
			}

			mc.ServerGroups[groupName].Mods = append(mc.ServerGroups[groupName].Mods, mod)
			err = ServerCfgSaver.Save()
		} else {
			InstallConfig.ClientMods = append(InstallConfig.ClientMods, mod)
			err = ConfigIo.Save(InstallConfig)
		}

		if err == nil {
			printToUser("Config updated.")
		}

		return
	},
}

func InitPrompts() {
	FriendlyPrompt = input.NewLinePrompt(AddFriendlyPromptText, input.NoOpValidator{})

	cliNameRegexValidator := input.NewRegexValidator(`^[a-z]+[a-z-]*[a-z]+$`,
		"must be two or more lowercase letters a-z; can include hyphens in between")
	cliNameUniqueValidator := &input.CliNameUniquenessValidator{GetModMap: getModMap}
	CliNamePrompt = input.NewLinePrompt(AddCliNamePromptText, cliNameRegexValidator, cliNameUniqueValidator)

	DescPrompt = input.NewLinePrompt(AddDescPromptText, input.NoOpValidator{})

	DetailsUrlPrompt = input.NewLinePrompt(AddDetailUrlPromptText, &input.UrlValidator{})

	DownloadUrlPrompt = input.NewLinePrompt(AddDownloadPromptText, &input.UrlValidator{})

	GroupPrompt = input.NewLinePrompt(AddGroupNamePromptText, &input.GroupNameValidator{})
}

// for testing
func ResetAddVars() {
	*serverMod = false
}

func init() {
	RootCmd.AddCommand(addCmd)

	flags := addCmd.Flags()

	serverMod = newBoolPtr(false)

	flags.BoolVar(serverMod, "server", false, "Add a new mod to the server config. Only allowed when building new versions of this tool.")

	InitPrompts()
}

func getModMap() mc.ModMap {
	return NameMapper.MapAllMods(InstallConfig.ClientMods)
}

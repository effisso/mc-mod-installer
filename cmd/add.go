package cmd

import (
	"mcmods/input"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

const (
	addFriendlyPromptText  = "What's the mod's name?\n> "
	addCliNamePromptText   = "Globally unique name (lowercase letters and hyphens only; short yet descriptive)\n> "
	addDescPromptText      = "Description of the mod (optional)\n> "
	addDetailURLPromptText = "Mod homepage/wiki URL\n> "
	addDownloadPromptText  = "Desired package download URL\n> "
	addGroupNamePromptText = "Server group\n> "
)

var (
	// FriendlyPrompt asks the user for a human-friendly name for the mod
	FriendlyPrompt input.Prompt

	// CliNamePrompt asks the user for a short, concise hyphenated name to
	// refer to the mod in CLI commands
	CliNamePrompt input.Prompt

	// DescPrompt prompts the user for a description of the mod
	DescPrompt input.Prompt

	// DetailsURLPrompt asks the user for a link to the mods homepage/wiki
	DetailsURLPrompt input.Prompt

	// DownloadURLPrompt asks the user for a link to the location where the
	// current version of the mod should be downloaded from
	DownloadURLPrompt input.Prompt

	// GroupPrompt prompts the user which group the mod should go in when
	// adding a new server mod
	GroupPrompt input.Prompt

	// ServerCfgSaver saves the server config (for building new versions of this tool)
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
		var friendlyName, cliName, desc, detURL, dlURL, groupName string

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

		if detURL, err = DetailsURLPrompt.GetInput(out, in); err != nil {
			return
		}

		if dlURL, err = DownloadURLPrompt.GetInput(out, in); err != nil {
			return
		}

		mod := &mc.Mod{
			FriendlyName: friendlyName,
			CliName:      cliName,
			Description:  desc,
			DetailsURL:   detURL,
			LatestURL:    dlURL,
		}

		if *serverMod {
			if groupName, err = GroupPrompt.GetInput(out, in); err != nil {
				return
			}

			mc.ServerGroups[groupName].Mods = append(mc.ServerGroups[groupName].Mods, mod)
			err = ServerCfgSaver.Save()
		} else {
			UserModConfig.ClientMods = append(UserModConfig.ClientMods, mod)
			err = cfgIo.Save(UserModConfig)
		}

		if err == nil {
			printToUser("Config updated.")
		}

		return
	},
}

// InitAddPrompts initializes all the prompts for the add command
func InitAddPrompts() {
	FriendlyPrompt = input.NewLinePrompt(addFriendlyPromptText, input.NoOpValidator{})

	cliNameRegexValidator := input.NewRegexValidator(`^[a-z]+[a-z-]*[a-z]+$`,
		"must be two or more lowercase letters a-z; can include hyphens in between")
	cliNameUniqueValidator := &input.CliNameUniquenessValidator{GetModMap: getModMap}
	CliNamePrompt = input.NewLinePrompt(addCliNamePromptText, cliNameRegexValidator, cliNameUniqueValidator)

	DescPrompt = input.NewLinePrompt(addDescPromptText, input.NoOpValidator{})

	DetailsURLPrompt = input.NewLinePrompt(addDetailURLPromptText, &input.URLValidator{})

	DownloadURLPrompt = input.NewLinePrompt(addDownloadPromptText, &input.URLValidator{})

	GroupPrompt = input.NewLinePrompt(addGroupNamePromptText, &input.GroupNameValidator{})
}

func init() {
	RootCmd.AddCommand(addCmd)

	flags := addCmd.Flags()

	serverMod = flags.Bool("server", false, "Add a new mod to the server config. Only allowed when building new versions of this tool.")

	InitAddPrompts()
}

func getModMap() mc.ModMap {
	return NameMapper.MapAllMods(UserModConfig.ClientMods)
}

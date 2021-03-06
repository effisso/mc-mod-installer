package cmd

import (
	"fmt"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	describeMap = map[string]func(string) error{}
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe [mod|group|install(ation)] <name>",
	Short: "Describe mods, mod installations, and server groups",
	Long: `
Describe prints out detailed information about the resource specified.
Examples:
 $ describe mod mod-name
 $ describe group required
 $ describe install mod-name`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rscType := args[0]
		name := args[1]

		describe, exists := describeMap[rscType]

		if exists {
			err := describe(name)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unknown resource: %s", rscType)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(describeCmd)

	describeMap["mod"] = describeMod

	describeMap["group"] = describeGroup

	describeMap["install"] = describeInstall
	describeMap["installation"] = describeInstall
}

func describeMod(modName string) error {
	cliMods := NameMapper.MapAllMods(UserModConfig.ClientMods)

	m := cliMods[modName]
	if m == nil {
		return mc.NewUnknownModError(modName)
	}

	printToUser(fmt.Sprintf("\n%s (%s)\n-----\n%s\nWebsite:  %s\nLatest package:  %s",
		m.FriendlyName, m.CliName, m.Description, m.DetailsURL, m.LatestURL))

	return nil
}

func describeGroup(groupName string) error {
	group := mc.ServerGroups[groupName]
	if group == nil {
		return mc.NewUnknownGroupError(groupName)
	}

	max := len(group.Mods) - 1
	for i, mod := range group.Mods {
		if i == max {
			printToUser(mod.CliName)
		} else {
			printLineToUser(mod.CliName)
		}
	}

	return nil
}

func describeInstall(modName string) error {
	cliMods := NameMapper.MapAllMods(UserModConfig.ClientMods)

	m := cliMods[modName]
	if m == nil {
		return mc.NewUnknownModError(modName)
	}

	i, exists := UserModConfig.ModInstallations[modName]

	if exists {
		printToUser(fmt.Sprintf("\n%s (%s)\n-----\nInstall timestamp:  %s\nUp-to-date:  %t",
			m.FriendlyName, m.CliName, i.Timestamp, m.LatestURL == i.DownloadURL))
	} else {
		printToUser("Not Installed.")
	}

	return nil
}

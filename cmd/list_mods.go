package cmd

import (
	"errors"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	listInstalled    = newBoolPtr(false)
	listNotInstalled = newBoolPtr(false)
	listClient       = newBoolPtr(false)
	listServer       = newBoolPtr(false)
	listGroup        = newStrPtr()
)

// modCmd represents the mod command
var modCmd = &cobra.Command{
	Use:   "mods",
	Short: "List mods by CLI name",
	Long: `
Prints mods by their unique CLI name. Filter by server groups, client/server
only, and installed or not.

Examples:
 $ list mods --installed
 $ list mods --client --not-installed
 $ list mods --group performance
 $ list mods --server --installed
 
 Providing both --installed and --not-installed is the same as providing
 neither. The --client and --server flags work similarly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !*listInstalled && !*listNotInstalled {
			*listInstalled = true
			*listNotInstalled = true
		}

		if *listGroup != "" {
			if *listClient {
				return errors.New("Can't specify a group with the client switch")
			}
			*listServer = true
			err := NameValidator.ValidateServerGroups([]string{*listGroup})
			if err != nil {
				return err
			}
		}

		if !*listClient && !*listServer {
			*listClient = true
			*listServer = true
		}

		if *listServer {
			for _, mod := range getServerMods() {
				if mod == nil {
					continue
				}
				printToUser(mod.CliName)
			}
		}

		if *listClient {
			for _, mod := range getClientMods() {
				printToUser(mod.CliName)
			}
		}
		return nil
	},
}

func init() {
	listCmd.AddCommand(modCmd)

	flags := modCmd.Flags()

	flags.BoolVarP(listInstalled, "installed", "i", false, "Show only mods that are installed currently.")

	flags.BoolVarP(listNotInstalled, "not-installed", "n", false, "Show only mods that are not installed currently.")

	flags.BoolVarP(listClient, "client", "c", false, "Show only client mods.")

	flags.BoolVarP(listServer, "server", "s", false, "Show only server mods.")

	flags.StringVarP(listGroup, "group", "g", "", "Show only mods from the specified group.")
}

// for testing
func ResetListVars() {
	*listInstalled = false
	*listNotInstalled = false
	*listClient = false
	*listServer = false
	*listGroup = ""
}

func getClientMods() []*mc.Mod {
	return getMods(InstallConfig.ClientMods, []*mc.Mod{})
}

func getServerMods() []*mc.Mod {
	mods := []*mc.Mod{}
	groupDefined := *listGroup != ""
	for groupName, group := range mc.ServerGroups {
		if groupDefined && *listGroup != groupName {
			continue
		}
		mods = getMods(group.Mods, mods)
	}
	return mods
}

func getMods(mods []*mc.Mod, apndTgt []*mc.Mod) []*mc.Mod {
	for _, mod := range mods {
		_, installed := InstallConfig.ModInstallations[mod.CliName]
		if *listInstalled && installed || *listNotInstalled && !installed {

			apndTgt = append(apndTgt, mod)
		}
	}

	return apndTgt
}

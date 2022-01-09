package cmd

import (
	"errors"
	"mcmods/mc"

	"github.com/spf13/cobra"
)

var (
	listInstalled    *bool
	listNotInstalled *bool
	listClient       *bool
	listServer       *bool
	listGroup        *string
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
				return errors.New("Can't specify a server group with the client switch")
			}
			*listServer = true
			if _, ok := mc.ServerGroups[*listGroup]; !ok {
				return mc.NewUnknownGroupError(*listGroup)
			}
		}

		if !*listClient && !*listServer {
			*listClient = true
			*listServer = true
		}

		if *listServer {
			for _, mod := range getServerMods() {
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

	listInstalled = flags.BoolP("installed", "i", false, "Show only mods that are installed currently.")

	listNotInstalled = flags.BoolP("not-installed", "n", false, "Show only mods that are not installed currently.")

	listClient = flags.BoolP("client", "c", false, "Show only client mods.")

	listServer = flags.BoolP("server", "s", false, "Show only server mods.")

	listGroup = flags.StringP("group", "g", "", "Show only mods from the specified group.")
}

// ResetListVars is used for testing
func ResetListVars() {
	*listInstalled = false
	*listNotInstalled = false
	*listClient = false
	*listServer = false
	*listGroup = ""
}

func getClientMods() []*mc.Mod {
	return getMods(UserModConfig.ClientMods, []*mc.Mod{})
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
		_, installed := UserModConfig.ModInstallations[mod.CliName]
		if *listInstalled && installed || *listNotInstalled && !installed {

			apndTgt = append(apndTgt, mod)
		}
	}

	return apndTgt
}

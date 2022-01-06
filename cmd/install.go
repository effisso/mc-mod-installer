package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

const (
	ServerOnlyGroupKey = "server-only"
)

var (
	CreateDownloaderFunc func(fs mc.FileSystem) mc.ModDownloader = mc.NewModDownloader

	NameValidator = mc.NewNameValidator()
	NameMapper    = mc.NewModNameMapper()
	Filter        = mc.NewModFilter(NameMapper, NameValidator)
	Installer     = mc.NewModInstaller()

	force      *bool
	fullServer *bool
	clientOnly *bool
	xMods      *[]string
	xGroups    *[]string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs and updates mod packages on the local machine.",
	Long: `
Install is used to initialize and update mod installations with this tool.

Users unfamiliar with configuring mods are encouraged to do a plain install
with all the required and recommended mods. Simply run the install command with
no arguments.

For advanced users wanting to use a custom set of performance-related mods or
those who simply don't want the optional mods on their machine, see the --help
command to read about filtering.

To perform a server install, use the --full-server option with nothing else:
$ install --full-server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !*fullServer {
			if *clientOnly {
				*xGroups = getServerModGroupNames(mc.ServerGroups)
			} else {
				*xGroups = append(*xGroups, ServerOnlyGroupKey)
			}
		}

		mods, err := Filter.FilterAllMods(*xGroups, *xMods, UserModConfig, *force)
		if err != nil {
			return err
		}

		dl := CreateDownloaderFunc(fs)
		err = Installer.InstallMods(dl, mods, UserModConfig)
		if err != nil {
			return err
		}

		err = cfgIo.Save(UserModConfig)
		if err != nil {
			return err
		}

		printToUser("Install completed.")
		return nil
	},
}

// ResetInstallVars is for testing
func ResetInstallVars() {
	*force = false
	*fullServer = false
	*clientOnly = false
	*xMods = (*xMods)[:0]
	*xGroups = (*xGroups)[:0]
}

func init() {
	RootCmd.AddCommand(installCmd)

	flags := installCmd.Flags()

	clientOnly = flags.BoolP("client-only", "c", false, "Only install your client mods.")

	force = flags.BoolP("force", "f", false, "Force the download and install, even if the file exists.")

	fullServer = flags.Bool("full-server", false, "Install all the server mods. Only necessary for the server itself. Ignores exclude flags.")

	xMods = flags.StringSlice("x-mod", []string{}, "Exclude specific Mods from the install (client or server). Specify multiple mods by separating the names with commas, no spaces.")

	xGroups = flags.StringSliceP("x-group", "x", []string{}, "Exclude Server Mod Groups from the install. The 'server-only' group is automatically excluded. Specify multiple mods by separating the names with commas, no spaces.")
}

func getServerModGroupNames(m map[string]*mc.ServerGroup) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

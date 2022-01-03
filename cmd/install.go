package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

const (
	ServerOnlyGroupKey = "server-only"
)

var (
	NameValidator = mc.NewNameValidator()
	NameMapper    = mc.NewModNameMapper()
	Filter        = mc.NewModFilter(NameMapper, NameValidator)
	Installer     = mc.NewModInstaller()
	Downloader    = mc.NewModDownloader()

	force      = newBoolPtr(false)
	fullServer = newBoolPtr(false)
	clientOnly = newBoolPtr(false)
	xMods      = []string{}
	xGroups    = []string{}
)

// installCmd represents the install command
var installCmd = NewInstallCmd()

func NewInstallCmd() *cobra.Command {
	return &cobra.Command{
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
					xGroups = getServerModGroupNames(mc.ServerGroups)
				} else {
					xGroups = append(xGroups, ServerOnlyGroupKey)
				}
			}

			mods, err := Filter.FilterAllMods(xGroups, xMods, InstallConfig, *force)
			if err != nil {
				return err
			}

			err = Installer.InstallMods(Downloader, mods, InstallConfig)
			if err != nil {
				return err
			}

			err = ConfigIo.Save(InstallConfig)
			if err != nil {
				return err
			}

			printToUser("Install completed.")
			return nil
		},
	}
}

// for testing
func ResetInstallVars() {
	*force = false
	*fullServer = false
	*clientOnly = false
	xMods = xMods[:0]
	xGroups = xGroups[:0]
}

func init() {
	RootCmd.AddCommand(installCmd)

	flags := installCmd.Flags()

	flags.BoolVarP(clientOnly, "client-only", "c", false, "Only install your client mods.")

	flags.BoolVarP(force, "force", "f", false, "Force the download and install, even if the file exists.")

	flags.BoolVar(fullServer, "full-server", false, "Install all the server mods. Only necessary for the server itself. Ignores exclude flags.")

	flags.StringSliceVar(&xMods, "x-mod", []string{}, "Exclude specific Mods from the install (client or server). Specify multiple mods by separating the names with commas, no spaces.")

	flags.StringSliceVarP(&xGroups, "x-group", "x", []string{}, "Exclude Server Mod Groups from the install. The 'server-only' group is automatically excluded. Specify multiple mods by separating the names with commas, no spaces.")
}

func getServerModGroupNames(m map[string]*mc.ServerGroup) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func newBoolPtr(val bool) *bool {
	b := val
	return &b
}

func newStrPtr() *string {
	s := ""
	return &s
}

package cmd

import (
	"mcmods/mc"

	"github.com/spf13/cobra"
)

const (
	// ServerOnlyGroupKey is the name of the mod group for mods which should only
	// be installed on the server
	ServerOnlyGroupKey = "server-only"
)

var (
	// CreateDownloaderFunc initializes the ModDownloader
	CreateDownloaderFunc func(fs mc.FileSystem) mc.ModDownloader = CreateDefaultDownloader

	// NameMapper creates a map of all mods to their CLI name
	NameMapper = mc.NewModNameMapper()

	// Filter returns a subset of mods based on input parameters
	Filter = mc.NewModFilter(NameMapper)

	// Installer installs mods
	Installer = mc.NewModInstaller()

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

IMPORTANT: Disconnect from VPN to avoid issues when downloading mods.

For advanced users wanting to use a custom set of performance-related mods or
those who simply don't want the optional mods on their machine, see the argument
descriptions for more about filtering out mods.

--force can be used to invoke a download even if the latest version of the mods
already exist locally. Otherwise, the tool skips if the latest URL matches the
URL at the time of download.

To perform a server install, use the --full-server option the FTP info:
  $ install --full-server --user <ftp-user> --password <pw> --server <server>

The FTP server and user are stored, so they're only needed on the first
command. Password is needed every time. `,
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

// CreateDefaultDownloader initializes a new mod downloader with a real HTTP Client
func CreateDefaultDownloader(fs mc.FileSystem) mc.ModDownloader {
	return mc.NewModDownloader(mc.NewHTTPClient(), fs)
}

package cmd

import (
	"fmt"
	"mcmods/mc"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	// CreateFsFunc creates the file system for FTP/Local. Exported for testing
	CreateFsFunc func(ftpArgs *mc.FtpArgs) (mc.FileSystem, error) = mc.NewFs

	// ConfigIoFunc instantiates an interface for config file IO
	ConfigIoFunc func(fs mc.FileSystem) mc.ModConfigIo = mc.NewUserModConfigIo

	// UserModConfig contains information about mod installations on the file
	// system
	UserModConfig *mc.UserModConfig

	// ViperInstance is the common instance of viper shared through the package
	ViperInstance = viper.GetViper()

	fs           mc.FileSystem
	cfgIo        mc.ModConfigIo
	cfgFile      string
	ftpUser      string
	ftpPw        string
	ftpServer    string
	printVersion *bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcmods",
	Short: "Tool for installing/maintaining mods for the CDP YAMS server.",
	Long: `
This tool installs and updates mods on a machine for connecting to the CDP YAMS
Minecraft server. The server is private, and only available by invite. To
inquire about an invite, please call 1-888-PISS-OFF and ask for Dianne.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var ftpArgs *mc.FtpArgs
		var err error

		if ftpPw != "" {
			ftpArgs = &mc.FtpArgs{
				Server: ViperInstance.GetString(mc.FtpServerKey),
				User:   ViperInstance.GetString(mc.FtpUserKey),
				Pw:     ftpPw,
			}
		}

		fs, err = CreateFsFunc(ftpArgs)
		cobra.CheckErr(err)

		cfgIo = ConfigIoFunc(fs)

		UserModConfig, err = cfgIo.LoadOrNew()
		cobra.CheckErr(err)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fs.Close()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	ResetRootVars()

	cobra.OnInitialize(initViper)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcmods.yaml)")

	RootCmd.PersistentFlags().StringVar(&ftpServer, "ftp-server", "", "The FTP server for managing server-side mods. Stored, only needed on the first command.")
	RootCmd.PersistentFlags().StringVarP(&ftpUser, "user", "u", "", "The FTP username. Stored, only needed on the first command.")
	RootCmd.PersistentFlags().StringVarP(&ftpPw, "password", "p", "", "The FTP password. Not stored, needed every time.")

	printVersion = RootCmd.Flags().BoolP("version", "v", false, "Show the version of this tool.")
}

// ResetRootVars is only used for testing
func ResetRootVars() {
	cfgFile = ""
	ftpUser = ""
	ftpPw = ""
	ftpServer = ""
}

// initViper reads in a config file through Viper
func initViper() {
	if cfgFile != "" {
		ViperInstance.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		ViperInstance.AddConfigPath(home)
		ViperInstance.SetConfigType("yaml")
		ViperInstance.SetConfigName(".mcmods")
	}

	ViperInstance.SetDefault(mc.InstallPathKey, mc.DefaultOsMinecraftDir)

	if err := ViperInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			err = ViperInstance.SafeWriteConfig()
			cobra.CheckErr(err)
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(err)
		}
	}

	updated := false
	if ftpUser != "" {
		updated = true
		ViperInstance.Set(mc.FtpUserKey, ftpUser)
	}
	if ftpServer != "" {
		updated = true
		ViperInstance.Set(mc.FtpServerKey, ftpServer)
	}

	if updated {
		ViperInstance.WriteConfig()
	}
}

// Print the complete output of the command to a user. Does not append a new
// line
func printToUser(txt string) {
	fmt.Fprint(RootCmd.OutOrStdout(), txt)
}

// Print a line of output for the command to a user. Appends a new line.
// Commands should conclude by calling printToUser to not append an empty
// final empty line to the buffer
func printLineToUser(txt string) {
	printToUser(fmt.Sprintf("%s\n", txt))
}

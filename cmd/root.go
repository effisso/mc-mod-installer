package cmd

import (
	"errors"
	"fmt"
	"mcmods/mc"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	ConfigIo = mc.NewModConfigIo()

	cfgFile       string
	InstallConfig *mc.ClientModConfig
	serverName    = "[SERVER]"
	ViperInstance = viper.GetViper()

	Version      string
	printVersion *bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcmods",
	Short: "Tool for installing/maintaining mods for the server " + serverName,
	Long: `
This tool installs and updates mods on a machine for connecting to.
a Minecraft server. The server is private, and only available by
invite. To inquire about an invite, please call 1-888-PISS-OFF and
ask for Dianne.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if *printVersion {
			fmt.Println(Version)
			return nil
		}
		return errors.New("no arguments specified")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcmods.yaml)")

	printVersion = RootCmd.Flags().BoolP("version", "v", false, "Show the version of this tool")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		ViperInstance.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mcmods" (without extension).
		ViperInstance.AddConfigPath(home)
		ViperInstance.SetConfigType("yaml")
		ViperInstance.SetConfigName(".mcmods")
	}

	ViperInstance.SetDefault(mc.InstallPathKey, mc.MinecraftDir)

	// If a config file is found, read it in.
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

	var err error
	InstallConfig, err = ConfigIo.LoadOrNew()
	if err != nil {
		panic(err)
	}
}

func printToUser(txt string) {
	fmt.Fprintln(RootCmd.OutOrStdout(), txt)
}

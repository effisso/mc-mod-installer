/*
Copyright © 2022 Zach Maddox <effisso@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"mcmods/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var serverName = "EffissoLand"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcmods",
	Short: "Tool for maintaining mods allowing connections to " + serverName,
	Long: "This tool installs and updates mods on a machine for connecting to " + serverName + `.
The server is private, and only available by invite. To inquire about an invite,
please call 1-888-PISS-OFF and ask for Dianne.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcmods.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	vpr := viper.GetViper()

	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mcmods" (without extension).
		vpr.AddConfigPath(home)
		vpr.SetConfigType("yaml")
		vpr.SetConfigName(".mcmods")
	}

	vpr.AutomaticEnv() // read in environment variables that match

	setDefaults()

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			err = vpr.SafeWriteConfig()
			cobra.CheckErr(err)
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(err)
		}
	}
}

func setDefaults() {
	vpr := viper.GetViper()

	var usrcfgloc, err = os.UserConfigDir()
	var installPath = filepath.Join(usrcfgloc, ".minecraft")

	if err == nil {
		vpr.SetDefault(config.McInstallPathKey, installPath)
	}
}

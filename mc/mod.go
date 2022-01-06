package mc

import (
	// embed needed for hard-coding server mod config into the tool
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/spf13/viper"
)

const (
	// ModsFolderName - The name of the mods folder in the minecraft installation directory
	ModsFolderName = "mods"

	// InstallPathKey - The key in the Viper config which defines the full path to Minecraft on disk
	InstallPathKey = "mcInstallPath"

	// FtpUserKey - The key of the FTP username
	FtpUserKey = "ftpUser"

	// FtpServerKey - The key of the FTP username
	FtpServerKey = "ftpServer"
)

var (
	//go:embed server_mods.json
	serverModJSON string

	// ServerGroups segregates mod definitions for the current version of this tool
	ServerGroups = map[string]*ServerGroup{}

	// ViperInstance - Shared instance of Viper for accessing config
	ViperInstance = viper.GetViper()
)

func init() {
	err := json.Unmarshal([]byte(serverModJSON), &ServerGroups)
	if err != nil {
		panic(errors.New("server_mods.json file couldn't be unmarshalled"))
	}
}

// Mod is a single downloadable JAR file representing a Minecraft mod
type Mod struct {
	FriendlyName string `json:"friendlyName"`
	CliName      string `json:"cliName"`
	Description  string `json:"description"`
	DetailsURL   string `json:"detailsUrl"`
	LatestURL    string `json:"latestUrl"`
}

// ServerGroup is a logical grouping of Mods on the Server
type ServerGroup struct {
	Description string `json:"description"`
	Mods        []*Mod `json:"mods"`
}

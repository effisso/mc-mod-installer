package mc

import (
	// embed needed for hard-coding server mod config into the tool
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	// ModFolderName - The name of the mods folder in the minecraft installation directory
	ModFolderName = "mods"

	// InstallPathKey - The key in the Viper config which defines the full path to Minecraft on disk
	InstallPathKey = "mcInstallPath"

	// FTPUserKey - The key of the FTP username
	FTPUserKey = "ftpUser"

	// FTPServerKey - The key of the FTP username
	FTPServerKey = "ftpServer"
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

// NewUnknownModError creates a new error indicating that the mod name provided
// by the user is not valid.
func NewUnknownModError(name string) error {
	return fmt.Errorf("Unknown Mod: %s", name)
}

// NewUnknownGroupError creates a new error indicating that the group name
// provided by the user is not valid.
func NewUnknownGroupError(name string) error {
	return fmt.Errorf("Unknown Server Group: %s", name)
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

package mc

import (
	_ "embed"
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// The name of the mods folder in the minecraft installation directory
	ModsFolderName = "mods"

	// The key in the Viper config which defines the full path to Minecraft on disk
	InstallPathKey = "mcInstallPath"
)

var (
	//go:embed "server_mods.json"
	ServerModJson string

	// Server JARs defined for the current version of this tool, by group
	ServerGroups = map[string]*ServerGroup{}

	// Shared instance of Viper for accessing config
	ViperInstance = viper.GetViper()
)

func init() {
	err := json.Unmarshal([]byte(ServerModJson), &ServerGroups)
	if err != nil {
		panic(errors.New("server_mods.json file couldn't be unmarshalled"))
	}
}

// Mod is a single downloadable JAR file representing a Minecraft mod
type Mod struct {
	FriendlyName string `json:"friendlyName"`
	CliName      string `json:"cliName"`
	Description  string `json:"description"`
	DetailsUrl   string `json:"detailsUrl"`
	LatestUrl    string `json:"latestUrl"`
}

// ServerGroup is a logical grouping of Mods on the Server
type ServerGroup struct {
	Description string `json:"description"`
	Mods        []*Mod `json:"mods"`
}

// RootDir returns the full path to the mods folder in the minecraft install. Install path is obtained from Viper (McInstallPathKey)
func RootDir() string {
	return filepath.Join(ViperInstance.GetString(InstallPathKey), ModsFolderName)
}

package mc

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	ModConfigFileName = "mcmods-install.json"
)

// ClientModConfig contains data about the mods installed by the tool on the system and custom client-only definitions of mods
type ClientModConfig struct {
	ModInstallations map[string]ModInstallation `json:"modInstallations"`
	ClientMods       []*Mod                     `json:"clientMods"`
}

// ModConfigIo interface for loading and saving the local installation config
type ModConfigIo interface {
	// LoadOrNew loads the JSON file and parses it, or returns a new config instance if not found
	LoadOrNew() (*ClientModConfig, error)

	// Save the config as JSON
	Save(cfg *ClientModConfig) error
}

type modConfigIo struct{}

func NewModConfigIo() ModConfigIo {
	return modConfigIo{}
}

func NewModConfig() ClientModConfig {
	return ClientModConfig{
		ModInstallations: map[string]ModInstallation{},
		ClientMods:       []*Mod{},
	}
}

// LoadOrNew loads the JSON file and parses it, or returns a new config instance if not found
func (m modConfigIo) LoadOrNew() (*ClientModConfig, error) {
	cfgPath := installConfigPath()
	cfg := &ClientModConfig{}
	bytes, err := afero.ReadFile(FileSystem, cfgPath)

	if err == nil {
		err = json.Unmarshal(bytes, cfg)
		if err != nil {
			cfg = nil
		}
	} else if os.IsNotExist(err) {
		cfg.ClientMods = []*Mod{}
		cfg.ModInstallations = map[string]ModInstallation{}
		err = nil
	}

	return cfg, err
}

// Save the config as JSON
func (m modConfigIo) Save(cfg *ClientModConfig) error {
	cfgPath := installConfigPath()
	bytes, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err // uncovered code :( don't know how to test this
	}

	return afero.WriteFile(FileSystem, cfgPath, bytes, 0644)
}

func installConfigPath() string {
	return filepath.Join(RootDir(), ModConfigFileName)
}

type ServerConfigSaver interface {
	Save() error
}

type serverConfigSaver struct{}

func (s serverConfigSaver) Save() error {
	return nil
}

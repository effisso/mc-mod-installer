package mc

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	// ModConfigFileName is the default name of the user config file
	ModConfigFileName = "mcmods-install.json"
)

// UserModConfig contains data about the mods installed by the tool on the
// system and custom client-only definitions of mods
type UserModConfig struct {
	ModInstallations map[string]ModInstallation `json:"modInstallations"`
	ClientMods       []*Mod                     `json:"clientMods"`
}

// ModConfigIo interface for loading and saving the local installation config
type ModConfigIo interface {
	// LoadOrNew loads the JSON file and parses it, or returns a new config
	// instance if not found
	LoadOrNew() (*UserModConfig, error)

	// Save the config as JSON
	Save(cfg *UserModConfig) error
}

type modConfigIo struct {
	Fs FileSystem
}

// NewUserModConfigIo returns a new interface for reading/writing mod config
func NewUserModConfigIo(fs FileSystem) ModConfigIo {
	return modConfigIo{Fs: fs}
}

// NewUserModConfig returns a new config with fields initialized and empty
func NewUserModConfig() UserModConfig {
	return UserModConfig{
		ModInstallations: map[string]ModInstallation{},
		ClientMods:       []*Mod{},
	}
}

// LoadOrNew loads the JSON file and parses it, or returns a new config instance if not found
func (m modConfigIo) LoadOrNew() (*UserModConfig, error) {
	cfg := &UserModConfig{}
	bytes, err := m.Fs.ReadFile(relUserConfigPath())

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
func (m modConfigIo) Save(cfg *UserModConfig) error {
	b, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	return m.Fs.WriteFile(bytes.NewReader(b), relUserConfigPath())
}

func relUserConfigPath() string {
	return filepath.Join(ModFolderName, ModConfigFileName)
}

// ServerConfigSaver saves the server mod JSON file which is embedded in this
// tool on build.
type ServerConfigSaver interface {
	Save() error
}

type serverConfigSaver struct{}

func (s serverConfigSaver) Save() error {
	return nil
}

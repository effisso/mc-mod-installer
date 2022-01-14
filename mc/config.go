package mc

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
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
		return err // not sure how to test :/
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

type serverConfigSaver struct {
	Afs afero.Fs
}

// NewServerConfigSaver creates a new implementation of the ServerConfigSaver
// interface with the given afero FS
func NewServerConfigSaver(fs afero.Fs) ServerConfigSaver {
	return &serverConfigSaver{Afs: fs}
}

func (s serverConfigSaver) Save() error {
	var b []byte
	var exists bool

	wd, err := os.Getwd()
	b, err = json.MarshalIndent(ServerGroups, "", "\t")
	if err != nil {
		return err // not sure how to test :/
	}

	cfgPath := filepath.Join(wd, "mc", "server_mods.json")

	exists, err = afero.Exists(s.Afs, cfgPath)
	if !exists {
		err = errors.New("adding server mods is only allowed when working from the root of the tool's code repo")
	}

	if err != nil {
		return err
	}

	return afero.WriteFile(s.Afs, cfgPath, b, 0644)
}

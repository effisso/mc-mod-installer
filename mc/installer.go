package mc

import (
	"fmt"
	"path/filepath"
	"time"
)

// ModInstallation captures the URL and filename for a mod that gets installed on the system
type ModInstallation struct {
	DownloadURL string `json:"downloadUrl"`
	Timestamp   string `json:"timestamp"`
}

// ModInstaller is an interface for for installing mods
type ModInstaller interface {
	// Downloads and installs the mods in the given slice
	InstallMods(downloader ModDownloader, mods []*Mod, cfg *UserModConfig) error
}

type modInstaller struct{}

// NewModInstaller returns a new struct which implements Installer
func NewModInstaller() ModInstaller {
	return modInstaller{}
}

// InstallMods downloads and installs the mods in the given slice
func (i modInstaller) InstallMods(downloader ModDownloader, mods []*Mod, cfg *UserModConfig) error {
	for _, m := range mods {
		fileName := fmt.Sprintf("%s.jar", m.CliName)
		modPath := filepath.Join(ModsFolderName, fileName)

		fmt.Printf("Installing %s\n", m.FriendlyName)
		err := downloader.Download(m, modPath)

		if err != nil {
			return err
		}

		cfg.ModInstallations[m.CliName] = ModInstallation{
			DownloadURL: m.LatestURL,
			Timestamp:   fmt.Sprint(time.Now().Format(time.UnixDate)),
		}
	}

	return nil
}

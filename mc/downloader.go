package mc

import (
	"fmt"
	"path/filepath"
)

// ModDownloader is the interface for downloading mods
type ModDownloader interface {
	// Download the mod to the specified file path
	Download(mod *Mod, filePath string) error
}

// ModDownloaderImpl only exported for testing access. Use ModDownloader interface
type ModDownloaderImpl struct {
	Fs         FileSystem
	HTTPClient *HTTPClient
}

// NewModDownloader creates a new instance of a struct which implements Downloader
// over the given http client
func NewModDownloader(hc *HTTPClient, fs FileSystem) ModDownloader {
	return &ModDownloaderImpl{
		Fs:         fs,
		HTTPClient: hc,
	}
}

// Download the specified mod from its LatestUrl and save it to the location specified
func (d ModDownloaderImpl) Download(mod *Mod, relPath string) error {
	err := d.Fs.MkDirAll(filepath.Dir(relPath))
	if err != nil {
		return err
	}

	fmt.Printf("  Downloading %s\n    to: %s\n", mod.LatestURL, relPath)

	resp, err := d.HTTPClient.Getter.Get(mod.LatestURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return d.Fs.WriteFile(resp.Body, relPath)
}

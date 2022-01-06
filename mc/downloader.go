package mc

import (
	"fmt"
	"net/http"
	"path/filepath"
)

var (
	// Http is the client for downloading mods
	Http = &HTTPClient{
		Getter: &http.Client{
			CheckRedirect: CheckRedirect,
		},
	}
)

// ModDownloader is the interface for downloading mods
type ModDownloader interface {
	// Download the mod to the specified file path
	Download(mod *Mod, filePath string) error
}

type modDownloader struct {
	fs FileSystem
}

// NewModDownloader creates a new instance of a struct which implements Downloader
func NewModDownloader(fs FileSystem) ModDownloader {
	return modDownloader{fs: fs}
}

// Download the specified mod from its LatestUrl and save it to the location specified
func (d modDownloader) Download(mod *Mod, relPath string) error {
	err := d.fs.MkDirAll(filepath.Dir(relPath))
	if err != nil {
		return err
	}

	fmt.Printf("  Downloading %s\n    to: %s\n", mod.LatestURL, relPath)

	resp, err := Http.Getter.Get(mod.LatestURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = d.fs.WriteFile(resp.Body, relPath)
	if err != nil {
		return err
	}

	return nil
}

// CheckRedirect makes redirects are followed. Only exported for testing
func CheckRedirect(r *http.Request, _ []*http.Request) error {
	r.URL.Opaque = r.URL.Path
	return nil
}

// HTTPGet is the interface for making HTTP GET requests abstractly
type HTTPGet interface {
	Get(url string) (*http.Response, error)
}

// HTTPClient contains interfaces for interacting with HTTP web services
type HTTPClient struct {
	Getter HTTPGet
}

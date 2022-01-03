package mc

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	// File system to use when reading/writing in the mc package
	FileSystem = afero.NewOsFs()
	Http       = &HttpClient{
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

type modDownloader struct{}

// NewDownloader creates a new instance of a struct which implements Downloader
func NewModDownloader() ModDownloader {
	return modDownloader{}
}

// Download the specified mod from its LatestUrl and save it to the location specified
func (d modDownloader) Download(mod *Mod, filePath string) error {
	err := FileSystem.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	fmt.Printf("  Downloading %s\n    to: %s\n", mod.LatestUrl, filePath)

	resp, err := Http.Getter.Get(mod.LatestUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = afero.WriteReader(FileSystem, filePath, resp.Body)
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

type HttpGet interface {
	Get(url string) (*http.Response, error)
}

type HttpClient struct {
	Getter HttpGet
}

func (h *HttpClient) GetRequest(url string) (*http.Response, error) {
	return h.Getter.Get(url)
}

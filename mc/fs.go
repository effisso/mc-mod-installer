package mc

import (
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// GetInstallPath reads the install path set in Viper.
func GetInstallPath() string {
	return ViperInstance.GetString(InstallPathKey)
}

// FileSystem is used to perform a few simple read and write operations for an
// arbitrary file system.
type FileSystem interface {
	WriteFile(r io.Reader, relPath string) error
	ReadFile(relPath string) ([]byte, error)
	MkDirAll(relPath string) error
	Close()
}

// LocalFileSystem reads and writes to the file system using afero.
type LocalFileSystem struct {
	Fs afero.Fs
}

// WriteFile writes the bytes to the relative path under the install directory.
// The file is given 0644 perms.
func (l LocalFileSystem) WriteFile(r io.Reader, relPath string) error {
	return afero.WriteReader(l.Fs, filepath.Join(GetInstallPath(), relPath), r)
}

// ReadFile reads the given relative path under the install directory.
func (l LocalFileSystem) ReadFile(relPath string) ([]byte, error) {
	return afero.ReadFile(l.Fs, filepath.Join(GetInstallPath(), relPath))
}

// MkDirAll creates all non-existant folders in the given path with 0755 perms
func (l LocalFileSystem) MkDirAll(relPath string) error {
	return l.Fs.MkdirAll(filepath.Join(GetInstallPath(), relPath), 0755)
}

// Close is a no-op for the local file system
func (l LocalFileSystem) Close() {}

// NewFs creates an FTPFileSystem if the args indicate FTP, or else
// a LocalFileSystem
func NewFs(ftpArgs *FTPArgs) (FileSystem, error) {
	var fs FileSystem
	if ftpArgs != nil {
		conn, err := openFTPToServer(ftpArgs)
		if err != nil {
			return nil, err
		}
		fs = &FTPFileSystem{Connection: conn}
	} else {
		fs = &LocalFileSystem{Fs: afero.NewOsFs()}
	}

	return fs, nil
}

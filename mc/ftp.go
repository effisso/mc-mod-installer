package mc

import (
	"errors"
	"io"
	"io/ioutil"
	"net/textproto"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

// FtpFileSystem is used to interact with Minecraft FTP servers to maintain mod
// installations.
type FtpFileSystem struct {
	Connection *ftp.ServerConn
}

// WriteFile writes the bytes over FTP to the given path on the server.
func (f FtpFileSystem) WriteFile(r io.Reader, relPath string) error {
	return f.Connection.Stor(fixPathForFtp(relPath), r)
}

// ReadFile reads the bytes of the given path over FTP.
func (f FtpFileSystem) ReadFile(relPath string) ([]byte, error) {
	protoErr := &textproto.Error{}
	r, err := f.Connection.Retr(fixPathForFtp(relPath))
	if err == nil {
		defer r.Close()
		return ioutil.ReadAll(r)
	} else if errors.As(err, &protoErr) {
		if protoErr.Code == ftp.StatusFileUnavailable {
			return nil, os.ErrNotExist
		}
	}
	return nil, err
}

// MkDirAll creates all non-existant folders in the given path with 0755 perms.
func (f FtpFileSystem) MkDirAll(relPath string) error {
	for _, dir := range GetRecursiveDirs(fixPathForFtp(relPath)) {
		if err := f.Connection.MakeDir(dir); err != nil {
			return err
		}
	}
	return nil
}

// Close calls Quit on the ftp connection
func (f FtpFileSystem) Close() {
	f.Connection.Quit()
}

// GetRecursiveDirs returns a slice of paths to each subdirectory in the dir
// hierarchy, in descending order. Exported only for testing purposes.
func GetRecursiveDirs(dir string) []string {
	if dir == "." || dir == "/" {
		return nil
	}
	dirs := []string{dir}
	for {
		if dir = path.Dir(dir); dir == "." || dir == "/" {
			break
		}
		dirs = append(dirs, dir)
	}

	// reverse for cleaner traversal
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}
	return dirs
}

// FtpArgs represents the information necessary to connect to FTP servers
type FtpArgs struct {
	Server string
	User   string
	Pw     string
}

func openFtpToServer(args *FtpArgs) (*ftp.ServerConn, error) {
	if args.Pw == "" || args.User == "" || args.Server == "" {
		return nil, errors.New("FTP access requires a username, password, and server")
	}

	ftpConnection, err := ftp.Dial(args.Server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = ftpConnection.Login(args.User, args.Pw)
	if err != nil {
		return nil, err
	}

	return ftpConnection, nil
}

func fixPathForFtp(path string) string {
	ps := string(os.PathSeparator)
	return "/" + strings.ReplaceAll(path, ps, "/")
}

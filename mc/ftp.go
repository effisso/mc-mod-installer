package mc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

var (
	// FTPDial is the function called when
	FTPDial func(server string, opts ...ftp.DialOption) (FTPConnection, error) = liveFTPDial
)

// FTPConnection is a wrapper for the ftp funcs used to send and receive FTP data
type FTPConnection interface {
	Login(username string, password string) error
	Stor(path string, r io.Reader) error
	Retr(path string) (*ftp.Response, error)
	MakeDir(dir string) error
	Quit() error
}

// FTPFileSystem is used to interact with Minecraft FTP servers to maintain mod
// installations.
type FTPFileSystem struct {
	Connection FTPConnection
}

// WriteFile writes the bytes over FTP to the given path on the server.
func (f FTPFileSystem) WriteFile(r io.Reader, relPath string) error {
	return f.Connection.Stor(fixPathForFTP(relPath), r)
}

// ReadFile reads the bytes of the given path over FTP.
func (f FTPFileSystem) ReadFile(relPath string) ([]byte, error) {
	protoErr := &textproto.Error{}
	r, err := f.Connection.Retr(fixPathForFTP(relPath))
	if err == nil {
		defer r.Close()
		return ioutil.ReadAll(r) // hard to unit test the happy path due to the library's architecture :(
	} else if errors.As(err, &protoErr) {
		if protoErr.Code == ftp.StatusFileUnavailable {
			return nil, os.ErrNotExist
		}
	}
	return nil, err
}

// MkDirAll creates all non-existant folders in the given path.
func (f FTPFileSystem) MkDirAll(relPath string) error {
	for _, dir := range GetRecursiveDirs(fixPathForFTP(relPath)) {
		if err := f.Connection.MakeDir(dir); err != nil {
			return err
		}
	}
	return nil
}

// Close calls Quit on the ftp connection
func (f FTPFileSystem) Close() {
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

// FTPArgs represents the information necessary to connect to FTP servers
type FTPArgs struct {
	Server    string
	User      string
	Pw        string
	TimeoutMs uint
}

func openFTPToServer(args *FTPArgs) (FTPConnection, error) {
	if args.Pw == "" || args.User == "" || args.Server == "" {
		return nil, errors.New("FTP access requires a username, password, and server")
	}

	fmt.Printf("Connecting FTP to %s\n", args.Server)

	timeoutOpt := ftp.DialWithTimeout(time.Duration(args.TimeoutMs) * time.Millisecond)

	ftpConnection, err := FTPDial(args.Server, timeoutOpt)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Logging in (%s)\n", args.User)

	err = ftpConnection.Login(args.User, args.Pw)
	if err != nil {
		return nil, err
	}

	return ftpConnection, nil
}

func fixPathForFTP(path string) string {
	ps := string(os.PathSeparator)
	return "/" + strings.ReplaceAll(path, ps, "/")
}

func liveFTPDial(server string, opts ...ftp.DialOption) (FTPConnection, error) {
	return ftp.Dial(server, opts...)
}

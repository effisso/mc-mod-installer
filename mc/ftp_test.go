package mc_test

import (
	"bytes"
	"errors"
	"io"
	"mcmods/mc"
	"net/textproto"
	"os"

	"github.com/jlaffaye/ftp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FTP File System", func() {
	Describe("NewFs", func() {
		var mock *mockFTP

		BeforeEach(func() {
			mock = emptyMock()
			mc.FTPDial = func(server string, opts ...ftp.DialOption) (mc.FTPConnection, error) {
				return mock, nil
			}
		})

		Context("missing FTP args", func() {
			It("returns an error when username is empty", func() {
				args := &mc.FTPArgs{Server: "server", Pw: "pw"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})

			It("returns an error when password is empty", func() {
				args := &mc.FTPArgs{Server: "server", User: "user"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})

			It("returns an error when server is empty", func() {
				args := &mc.FTPArgs{User: "user", Pw: "pw"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})
		})

		It("dials and logs into ftp if given args", func() {
			ftpfs := &mc.FTPFileSystem{}
			args := &mc.FTPArgs{
				Server:    "<server>",
				User:      "<usr>",
				Pw:        "<pw>",
				TimeoutMs: 123,
			}

			fs, err := mc.NewFs(args)

			Expect(fs).ToNot(BeNil())
			Expect(err).To(BeNil())
			Expect(fs).To(BeAssignableToTypeOf(ftpfs))
		})

		It("returns errors from ftp dial", func() {
			dialErr := errors.New("dial error")
			args := &mc.FTPArgs{
				Server:    "<server>",
				User:      "<usr>",
				Pw:        "<pw>",
				TimeoutMs: 123,
			}
			mc.FTPDial = func(server string, opts ...ftp.DialOption) (mc.FTPConnection, error) {
				return nil, dialErr
			}

			_, err := mc.NewFs(args)

			Expect(err).To(Equal(dialErr))
		})

		It("returns errors from ftp login", func() {
			loginErr := errors.New("login error")
			args := &mc.FTPArgs{
				Server:    "<server>",
				User:      "<usr>",
				Pw:        "<pw>",
				TimeoutMs: 123,
			}
			mock.LoginFunc = func(username, password string) error {
				return loginErr
			}

			_, err := mc.NewFs(args)

			Expect(err).To(Equal(loginErr))
		})
	})

	Describe("FTP File System", func() {
		var ftpFs *mc.FTPFileSystem
		var mock *mockFTP
		var r io.Reader

		BeforeEach(func() {
			mock = emptyMock()
			ftpFs = &mc.FTPFileSystem{
				Connection: mock,
			}

			r = bytes.NewReader([]byte("test file content"))
		})

		Context("WriteFile", func() {
			It("calls stor with the bytes and fixed path", func() {
				relPath := `some/path/to/a.file`
				called := false
				mock.StorFunc = func(path string, r io.Reader) error {
					called = true
					Expect(path).To(Equal("/" + relPath))
					return nil
				}

				err := ftpFs.WriteFile(r, relPath)

				Expect(err).To(BeNil())
				Expect(called).To(BeTrue())
			})

			It("returns errors from stor", func() {
				storErr := errors.New("stor error")
				relPath := "path/to/file.txt"
				mock.StorFunc = func(path string, r io.Reader) error {
					return storErr
				}

				err := ftpFs.WriteFile(r, relPath)

				Expect(err).To(Equal(storErr))
			})
		})

		Context("ReadFile", func() {
			It("returns errors from retr", func() {
				retrErr := errors.New("retr error")
				mock.RetrFunc = func(path string) (*ftp.Response, error) {
					return nil, retrErr
				}

				_, err := ftpFs.ReadFile("path/to/file.txt")

				Expect(err).To(Equal(retrErr))
			})

			It("converts FTP error 550 (file unavailable) into an os.ErrNotExist", func() {
				mock.RetrFunc = func(path string) (*ftp.Response, error) {
					return nil, &textproto.Error{Code: ftp.StatusFileUnavailable}
				}

				_, err := ftpFs.ReadFile("another/path/to/file.txt")

				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})

		Context("MkDirAll", func() {
			folderPath := "path/with/several/dirs"

			It("returns errors from the MakeDir call", func() {
				mkDirErr := errors.New("mkdir error")
				mock.MakeDirFunc = func(dir string) error {
					return mkDirErr
				}

				err := ftpFs.MkDirAll(folderPath)

				Expect(err).To(Equal(mkDirErr))
			})

			It("makes each level of the hierarchy", func() {
				curIndex := 0
				expectedDirsInOrder := []string{
					"/path",
					"/path/with",
					"/path/with/several",
					"/" + folderPath,
				}
				mock.MakeDirFunc = func(dir string) error {
					Expect(dir).To(Equal(expectedDirsInOrder[curIndex]))
					curIndex++
					return nil
				}

				err := ftpFs.MkDirAll(folderPath)

				Expect(err).To(BeNil())
				Expect(curIndex).To(Equal(len(expectedDirsInOrder)))
			})

			It("works for just the root directory", func() {
				called := false
				mock.MakeDirFunc = func(dir string) error {
					called = true
					Expect(dir).To(Equal("hello"))
					return nil
				}

				err := ftpFs.MkDirAll("")

				Expect(err).To(BeNil())
				Expect(called).To(BeFalse(), "MakeDir shoudn't be called on the root")
			})

			It("doesn't fail when the directory already exists", func() {
				mock.MakeDirFunc = func(dir string) error {
					return &textproto.Error{Code: ftp.StatusFileUnavailable}
				}

				err := ftpFs.MkDirAll("existing")

				Expect(err).To(BeNil())
			})
		})

		Context("Close", func() {
			It("calls Quit", func() {
				called := false
				mock.QuitFunc = func() error {
					called = true
					return nil
				}

				ftpFs.Close()

				Expect(called).To(BeTrue())
			})
		})
	})
})

type mockFTP struct {
	LoginFunc   func(username string, password string) error
	StorFunc    func(path string, r io.Reader) error
	RetrFunc    func(path string) (*ftp.Response, error)
	MakeDirFunc func(dir string) error
	QuitFunc    func() error
}

func emptyMock() *mockFTP {
	return &mockFTP{
		LoginFunc:   func(username string, password string) error { return nil },
		StorFunc:    func(path string, r io.Reader) error { return nil },
		RetrFunc:    func(path string) (*ftp.Response, error) { return &ftp.Response{}, nil },
		MakeDirFunc: func(dir string) error { return nil },
		QuitFunc:    func() error { return nil },
	}
}

func (ftp mockFTP) Login(username string, password string) error {
	return ftp.LoginFunc(username, password)
}

func (ftp mockFTP) Stor(path string, r io.Reader) error {
	return ftp.StorFunc(path, r)
}

func (ftp mockFTP) Retr(path string) (*ftp.Response, error) {
	return ftp.RetrFunc(path)
}

func (ftp mockFTP) MakeDir(dir string) error {
	return ftp.MakeDirFunc(dir)
}

func (ftp mockFTP) Quit() error {
	return ftp.QuitFunc()
}

func fakeFTPDial(server string, opts ...ftp.DialOption) (mc.FTPConnection, error) {
	return emptyMock(), nil
}

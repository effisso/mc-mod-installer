package mc_test

import (
	"bytes"
	"errors"
	"io"
	"mcmods/mc"
	"net/textproto"
	"os"
	"strings"

	"github.com/jlaffaye/ftp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FTP File System", func() {
	Describe("NewFs", func() {
		var mock *mockFtp

		BeforeEach(func() {
			mock = emptyMock()
			mc.FTPDial = func(server string, opts ...ftp.DialOption) (mc.FTPConnection, error) {
				return mock, nil
			}
		})

		Context("missing FTP args", func() {
			It("returns an error when username is empty", func() {
				args := &mc.FtpArgs{Server: "server", Pw: "pw"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})

			It("returns an error when password is empty", func() {
				args := &mc.FtpArgs{Server: "server", User: "user"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})

			It("returns an error when server is empty", func() {
				args := &mc.FtpArgs{User: "user", Pw: "pw"}

				_, err := mc.NewFs(args)

				Expect(err).ToNot(BeNil())
			})
		})

		It("dials and logs into ftp if given args", func() {
			ftpfs := &mc.FtpFileSystem{}
			args := &mc.FtpArgs{
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
			args := &mc.FtpArgs{
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
			args := &mc.FtpArgs{
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
		var ftpFs *mc.FtpFileSystem
		var mock *mockFtp
		var r io.Reader

		BeforeEach(func() {
			mock = emptyMock()
			ftpFs = &mc.FtpFileSystem{
				Connection: mock,
			}

			r = bytes.NewReader([]byte("test file content"))
		})

		Context("WriteFile", func() {
			It("calls stor with the bytes and fixed path", func() {
				relPath := "windows\\formatted\\path\\should\\be.converted"
				called := false
				mock.StorFunc = func(path string, r io.Reader) error {
					called = true
					Expect(path).To(Equal("/" + strings.ReplaceAll(relPath, "\\", "/")))
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

		// Context("GetRecursiveDirs", func() {
		// 	It("Returns path hierarchy strings in order", func() {
		// 		path := "/test/directory/hierarchy/test"
		// 		expected := []string{
		// 			"/test",
		// 			"/test/directory",
		// 			"/test/directory/hierarchy",
		// 			path,
		// 		}

		// 		dirs := mc.GetRecursiveDirs(path)

		// 		Expect(dirs).To(HaveLen(len(expected)))
		// 		for i, dir := range dirs {
		// 			Expect(dir).To(Equal(expected[i]))
		// 		}
		// 	})

		// 	It("returns nil for root directories", func() {
		// 		roots := []string{".", "/"}

		// 		for _, root := range roots {
		// 			Expect(mc.GetRecursiveDirs(root)).To(BeNil(), "didn't return nil for root: "+root)
		// 		}
		// 	})
		// })
	})
})

type mockFtp struct {
	LoginFunc   func(username string, password string) error
	StorFunc    func(path string, r io.Reader) error
	RetrFunc    func(path string) (*ftp.Response, error)
	MakeDirFunc func(dir string) error
	QuitFunc    func() error
}

func emptyMock() *mockFtp {
	return &mockFtp{
		LoginFunc:   func(username string, password string) error { return nil },
		StorFunc:    func(path string, r io.Reader) error { return nil },
		RetrFunc:    func(path string) (*ftp.Response, error) { return &ftp.Response{}, nil },
		MakeDirFunc: func(dir string) error { return nil },
		QuitFunc:    func() error { return nil },
	}
}

func (ftp mockFtp) Login(username string, password string) error {
	return ftp.LoginFunc(username, password)
}

func (ftp mockFtp) Stor(path string, r io.Reader) error {
	return ftp.StorFunc(path, r)
}

func (ftp mockFtp) Retr(path string) (*ftp.Response, error) {
	return ftp.RetrFunc(path)
}

func (ftp mockFtp) MakeDir(dir string) error {
	return ftp.MakeDirFunc(dir)
}

func (ftp mockFtp) Quit() error {
	return ftp.QuitFunc()
}

func fakeFTPDial(server string, opts ...ftp.DialOption) (mc.FTPConnection, error) {
	return emptyMock(), nil
}

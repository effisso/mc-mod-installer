package mc_test

import (
	"errors"
	"io"
	"mcmods/mc"
	. "mcmods/testdata"
	"net/http"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Downloader", func() {
	BeforeEach(func() {
		InitTestData()
	})

	Context("download func", func() {
		var fs afero.Fs
		var mcfs *mc.LocalFileSystem
		var dl mc.ModDownloader
		var rc io.ReadCloser
		var hc *mc.HTTPClient
		var eg *emptyGetter

		mcInstallPath := "/root/folder/.minecraft"
		relFilePath := "some/path/to/a.jar"
		fullPath := filepath.Join(mcInstallPath, relFilePath)
		content := "test"

		BeforeEach(func() {
			mc.ViperInstance.Set(mc.InstallPathKey, mcInstallPath)

			rc = io.NopCloser(strings.NewReader(content))
			fs = afero.NewMemMapFs()
			mcfs = &mc.LocalFileSystem{Fs: fs}
			eg = &emptyGetter{Res: &http.Response{Body: rc}}
			hc = &mc.HTTPClient{Getter: eg}
			dl = mc.NewModDownloader(hc, mcfs)
		})

		It("creates directories if not present, writes file contents", func() {
			err := dl.Download(TestingClientMod1, relFilePath)

			Expect(err).To(BeNil())

			exists, _ := afero.Exists(fs, fullPath)
			Expect(exists).To(BeTrue())

			b, _ := afero.ReadFile(fs, fullPath)
			Expect(string(b)).To(Equal(content))
		})

		It("doesn't download if dirs can't be created", func() {
			eg.Err = errors.New("this error won't be returned")
			mcfs.Fs = afero.NewReadOnlyFs(fs)

			err := dl.Download(TestingClientMod1, relFilePath)

			Expect(err).To(Not(BeNil()))
			Expect(err).To(Not(Equal(eg.Err)))

			exists, _ := afero.Exists(fs, fullPath)
			Expect(exists).To(BeFalse())
		})

		It("doesn't write the file if the download fails", func() {
			eg.Err = errors.New("bad url, or something. idk")

			err := dl.Download(TestingClientMod1, relFilePath)

			Expect(err).To(Equal(eg.Err))

			exists, _ := afero.Exists(fs, fullPath)
			Expect(exists).To(BeFalse())
		})

		It("returns an error if the write fails", func() {
			hc.Getter = emptyGetterWithTask{
				emptyGetter: emptyGetter{Res: &http.Response{Body: rc}},
				task: func() {
					// while the http get would be happening, "lock" the file system so the file
					// can't be written to ensure errors created from the underlying FS get returned
					mcfs.Fs = afero.NewReadOnlyFs(fs)
				},
			}

			err := dl.Download(TestingClientMod1, relFilePath)

			Expect(err).To(Not(BeNil()))

			exists, _ := afero.Exists(fs, fullPath)
			Expect(exists).To(BeFalse())
		})
	})
})

// -----
// FAKE DOWNLOADERS
// -----

// do nothing, optionally return an error
type emptyDownloader struct {
	Err error
}

func (e emptyDownloader) Download(mod *mc.Mod, modFolder string) error {
	return e.Err
}

// verify Download func args
type verifyingDownloader struct {
	ExpectedPath string
	ExpectedMod  *mc.Mod
}

func (v verifyingDownloader) Download(mod *mc.Mod, modFolder string) error {
	Expect(mod).To(Equal(v.ExpectedMod))
	Expect(modFolder).To(Equal(v.ExpectedPath))
	return nil
}

// count calls to the Download func
type countingDownloader struct {
	CallCount int
}

func (c *countingDownloader) Download(mod *mc.Mod, modFolder string) error {
	c.CallCount++
	return nil
}

// -----
// FAKE HTTP CLIENTS
// -----

type emptyGetter struct {
	Res *http.Response
	Err error
}

// just return the response and error on the struct
func (g emptyGetter) Get(url string) (*http.Response, error) {
	return g.Res, g.Err
}

// verify the url passed into the Get func
type getURLVerifier struct {
	emptyGetter
	ExpectedURL string
}

// check the url and return the vals on the struct
func (v getURLVerifier) Get(url string) (*http.Response, error) {
	Expect(url).To(Equal(v.ExpectedURL))
	return v.Res, v.Err
}

// runs a task when the download should happen to change some state during a test
type emptyGetterWithTask struct {
	emptyGetter
	task func()
}

// just return the response and error on the struct
func (g emptyGetterWithTask) Get(url string) (*http.Response, error) {
	g.task()
	return g.Res, g.Err
}

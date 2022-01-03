package mc_test

import (
	"errors"
	"io"
	"mcmods/mc"
	. "mcmods/testdata"
	"net/http"
	"net/url"
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Downloader", func() {
	BeforeEach(func() {
		InitTestData()
	})

	Context("http client", func() {
		It("follows redirects", func() {
			testPath := "http://www.google.com/"
			req := &http.Request{
				URL: &url.URL{
					Path: testPath,
				},
			}

			client, ok := mc.Http.Getter.(*http.Client)
			Expect(ok).To(BeTrue())

			err := client.CheckRedirect(req, []*http.Request{})

			Expect(err).To(BeNil())
			Expect(req.URL.Opaque).To(Equal(testPath))
		})
	})

	Context("download func", func() {
		var dl mc.ModDownloader
		var rc io.ReadCloser

		filePath := "/some/path/to/a.jar"
		content := "test"

		BeforeEach(func() {
			dl = mc.NewModDownloader()
			mc.FileSystem = afero.NewMemMapFs()
			rc = io.NopCloser(strings.NewReader(content))
		})

		It("creates directories if not present, writes file contents", func() {
			mc.Http.Getter = emptyGetter{Res: &http.Response{Body: rc}}

			err := dl.Download(TestingClientMod1, filePath)

			Expect(err).To(BeNil())

			exists, _ := afero.Exists(mc.FileSystem, filePath)
			Expect(exists).To(BeTrue())

			b, _ := afero.ReadFile(mc.FileSystem, filePath)
			Expect(string(b)).To(Equal(content))
		})

		It("doesn't download if dirs can't be created", func() {
			eg := emptyGetter{Err: errors.New("this error won't be returned")}
			mc.FileSystem = afero.NewReadOnlyFs(mc.FileSystem)
			mc.Http.Getter = eg

			err := dl.Download(TestingClientMod1, filePath)

			Expect(err).To(Not(BeNil()))
			Expect(err).To(Not(Equal(eg.Err)))

			exists, _ := afero.Exists(mc.FileSystem, filePath)
			Expect(exists).To(BeFalse())
		})

		It("doesn't write the file if the download fails", func() {
			eg := emptyGetter{Err: errors.New("bad url, or something. idk")}
			mc.Http.Getter = eg

			err := dl.Download(TestingClientMod1, filePath)

			Expect(err).To(Equal(eg.Err))

			exists, _ := afero.Exists(mc.FileSystem, filePath)
			Expect(exists).To(BeFalse())
		})

		It("returns an error if the write fails", func() {
			eg := emptyGetter{Res: &http.Response{Body: rc}}
			mc.Http.Getter = eg

			mc.FileSystem.MkdirAll(path.Dir(filePath), 0755)
			mc.FileSystem = afero.NewReadOnlyFs(mc.FileSystem)

			err := dl.Download(TestingClientMod1, filePath)

			Expect(err).To(Not(BeNil()))

			exists, _ := afero.Exists(mc.FileSystem, filePath)
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
type getUrlVerifier struct {
	emptyGetter
	ExpectedUrl string
}

// check the url and return the vals on the struct
func (v getUrlVerifier) Get(url string) (*http.Response, error) {
	Expect(url).To(Equal(v.ExpectedUrl))
	return v.Res, v.Err
}

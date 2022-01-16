package mc_test

import (
	"mcmods/mc"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP", func() {
	Context("http client", func() {
		It("follows redirects", func() {
			testPath := "http://www.google.com/"
			req := &http.Request{
				URL: &url.URL{
					Path: testPath,
				},
			}

			client := mc.NewHTTPClient()
			clientGetter, ok := client.Getter.(*http.Client)
			Expect(ok).To(BeTrue())

			err := clientGetter.CheckRedirect(req, []*http.Request{})

			Expect(err).To(BeNil())
			Expect(req.URL.Opaque).To(Equal(testPath))
		})

		It("converts + sign to %2B", func() {
			testPath := "http://www.website.com/my+package.jar"
			expectedPath := "http://www.website.com/my%2Bpackage.jar"
			req := &http.Request{
				URL: &url.URL{
					Path: testPath,
				},
			}

			client := mc.NewHTTPClient()
			clientGetter, ok := client.Getter.(*http.Client)
			Expect(ok).To(BeTrue())

			err := clientGetter.CheckRedirect(req, []*http.Request{})

			Expect(err).To(BeNil())
			Expect(req.URL.Opaque).To(Equal(expectedPath))
		})
	})
})

package mc_test

import (
	"mcmods/mc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FTP File System", func() {

	// ----------
	// Test coverage here is shoddy because it's not easy to mock the FTP library
	// ----------

	Describe("NewFs", func() {
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

		Context("dial errors", func() {
			It("returns dial errors", func() {

			})
		})
	})

	Describe("FTP File System", func() {
		Context("GetRecursiveDirs", func() {
			It("Returns path hierarchy strings in order", func() {
				path := "/test/directory/hierarchy/test"
				expected := []string{
					"/test",
					"/test/directory",
					"/test/directory/hierarchy",
					path,
				}

				dirs := mc.GetRecursiveDirs(path)

				Expect(dirs).To(HaveLen(len(expected)))
				for i, dir := range dirs {
					Expect(dir).To(Equal(expected[i]))
				}
			})

			It("returns nil for root directories", func() {
				roots := []string{".", "/"}

				for _, root := range roots {
					Expect(mc.GetRecursiveDirs(root)).To(BeNil(), "didn't return nil for root: "+root)
				}
			})
		})
	})
})

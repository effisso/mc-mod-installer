package mc_test

import (
	"bytes"
	"mcmods/mc"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Local File System", func() {
	var fs mc.FileSystem
	var aferoMemMap afero.Fs
	mcInstallLoc := "/test/path/minecraft"
	relPath := mc.ModFolderName + "/test.txt"
	fullPath := mcInstallLoc + "/" + relPath
	fileContent := "test file content"
	expectedBytes := []byte(fileContent)

	BeforeEach(func() {
		aferoMemMap = afero.NewMemMapFs()
		fs = &mc.LocalFileSystem{Fs: aferoMemMap}

		mc.ViperInstance.Set(mc.InstallPathKey, mcInstallLoc)
	})

	It("NewFs returns a new file system with no ftp args", func() {
		fs, err := mc.NewFs(nil)
		Expect(fs).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	Context("WriteFile", func() {
		It("builds abs path from relpath param", func() {
			Expect(aferoMemMap.MkdirAll(filepath.Dir(fullPath), 0755)).To(BeNil())

			Expect(fs.WriteFile(bytes.NewReader(expectedBytes), relPath)).To(BeNil())

			readBytes, _ := afero.ReadFile(aferoMemMap, fullPath)
			Expect(string(readBytes)).To(Equal(fileContent))
		})
	})

	Context("ReadFile", func() {
		It("builds abs path from relpath param", func() {
			Expect(aferoMemMap.MkdirAll(filepath.Dir(fullPath), 0755)).To(BeNil())
			Expect(afero.WriteFile(aferoMemMap, fullPath, expectedBytes, 0644)).To(BeNil())

			readBytes, err := fs.ReadFile(relPath)

			Expect(err).To(BeNil())
			Expect(string(readBytes)).To(Equal(fileContent))
		})
	})

	Context("MkDirAll", func() {
		It("builds abs path from relpath param", func() {
			Expect(fs.MkDirAll(mcInstallLoc)).To(BeNil())

			exists, err := afero.Exists(aferoMemMap, mcInstallLoc)

			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})
	})

	Context("Close", func() {
		It("does nothing", func() {
			fs.Close()
			Expect(1).To(Equal(1))
		})
	})
})

package mc_test

import (
	"mcmods/mc"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMods(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MC Suite")
}

var _ = Describe("Mod", func() {
	Describe("Root Path", func() {
		It("should be built from the mc install location in Viper", func() {
			pathValue := "/some/path"
			expected := filepath.Join(pathValue, mc.ModsFolderName)
			mc.ViperInstance.Set(mc.InstallPathKey, pathValue)

			Expect(mc.RootDir()).To(Equal(expected))
		})
	})
})

package mc_test

import (
	"mcmods/mc"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMods(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MC Suite")
}

var _ = Describe("Mod", func() {
	Describe("Install Path", func() {
		It("should be read from the mc install location in Viper", func() {
			pathValue := "/some/path"

			mc.ViperInstance.Set(mc.InstallPathKey, pathValue)

			Expect(mc.GetInstallPath()).To(Equal(pathValue))
		})
	})
})

package mc_test

import (
	"encoding/json"
	"mcmods/mc"
	. "mcmods/testdata"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Config IO", func() {
	var configIo mc.ModConfigIo
	var configPath, content, emptyContent string

	BeforeEach(func() {
		InitTestData()

		mcPath := "/path/to/mc"
		configPath = filepath.Join(mcPath, mc.ModsFolderName, mc.ModConfigFileName)

		b, _ := json.MarshalIndent(TestingConfig, "", "\t")
		content = string(b)

		b, _ = json.MarshalIndent(mc.ClientModConfig{}, "", "\t")
		emptyContent = string(b)

		configIo = mc.NewModConfigIo()
		mc.FileSystem = afero.NewMemMapFs()
		mc.ViperInstance.Set(mc.InstallPathKey, mcPath)
	})

	Context("LoadOrNew func", func() {
		When("the file exists", func() {
			When("the file has no installs/mods", func() {
				It("unmarshals from the file", func() {
					afero.WriteFile(mc.FileSystem, configPath, []byte(emptyContent), 0644)

					cfg, err := configIo.LoadOrNew()

					Expect(err).To(BeNil())
					Expect(cfg.ClientMods).To(BeEmpty())
					Expect(cfg.ModInstallations).To(BeEmpty())
				})
			})

			When("has installs/mods", func() {
				It("unmarshals from the file", func() {
					afero.WriteFile(mc.FileSystem, configPath, []byte(content), 0644)

					cfg, err := configIo.LoadOrNew()

					Expect(err).To(BeNil())
					Expect(cfg.ClientMods).To(HaveLen(len(TestingConfig.ClientMods)))

					for i := 0; i < len(cfg.ClientMods); i++ {
						cfgMod := cfg.ClientMods[i]
						xpctMod := TestingConfig.ClientMods[i]
						Expect(cfgMod.CliName).To(Equal(xpctMod.CliName))
						Expect(cfgMod.FriendlyName).To(Equal(xpctMod.FriendlyName))
						Expect(cfgMod.Description).To(Equal(xpctMod.Description))
						Expect(cfgMod.DetailsUrl).To(Equal(xpctMod.DetailsUrl))
						Expect(cfgMod.LatestUrl).To(Equal(xpctMod.LatestUrl))
					}

					Expect(cfg.ModInstallations).To(HaveLen(len(TestingConfig.ModInstallations)))
					for cliName, cfgMod := range cfg.ModInstallations {
						xpctMod := TestingConfig.ModInstallations[cliName]
						Expect(cfgMod.DownloadUrl).To(Equal(xpctMod.DownloadUrl))
						Expect(cfgMod.Timestamp).To(Equal(xpctMod.Timestamp))
					}
				})
			})
		})

		When("file does not exist", func() {
			It("returns an empty config", func() {
				cfg, err := configIo.LoadOrNew()

				Expect(err).To(BeNil())
				Expect(cfg.ClientMods).To(BeEmpty())
				Expect(cfg.ModInstallations).To(BeEmpty())
			})
		})

		When("there are json unmarshalling errors", func() {
			It("returns the error", func() {
				afero.WriteFile(mc.FileSystem, configPath, []byte("{"), 0644)

				cfg, err := configIo.LoadOrNew()

				Expect(cfg).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("Save func", func() {
		It("saves empty config", func() {
			err := configIo.Save(&mc.ClientModConfig{})

			Expect(err).To(BeNil())

			b, err := afero.ReadFile(mc.FileSystem, configPath)

			Expect(string(b)).To(Equal(emptyContent))
		})

		It("saves non-empty config", func() {
			err := configIo.Save(TestingConfig)

			Expect(err).To(BeNil())

			b, err := afero.ReadFile(mc.FileSystem, configPath)

			Expect(string(b)).To(Equal(content))
		})
	})
})

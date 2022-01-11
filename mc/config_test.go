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
	var fs afero.Fs

	BeforeEach(func() {
		InitTestData()

		mcPath := "/path/to/mc"
		configPath = filepath.Join(mcPath, mc.ModFolderName, mc.ModConfigFileName)

		b, _ := json.MarshalIndent(TestingConfig, "", "\t")
		content = string(b)

		b, _ = json.MarshalIndent(mc.UserModConfig{}, "", "\t")
		emptyContent = string(b)

		fs = afero.NewMemMapFs()
		mcfs := &mc.LocalFileSystem{Fs: fs}
		configIo = mc.NewUserModConfigIo(mcfs)
		mc.ViperInstance.Set(mc.InstallPathKey, mcPath)
	})

	Context("LoadOrNew func", func() {
		When("the file exists", func() {
			When("the file has no installs/mods", func() {
				It("unmarshals from the file", func() {
					afero.WriteFile(fs, configPath, []byte(emptyContent), 0644)

					cfg, err := configIo.LoadOrNew()

					Expect(err).To(BeNil())
					Expect(cfg.ClientMods).To(BeEmpty())
					Expect(cfg.ModInstallations).To(BeEmpty())
				})
			})

			When("has installs/mods", func() {
				It("unmarshals from the file", func() {
					afero.WriteFile(fs, configPath, []byte(content), 0644)

					cfg, err := configIo.LoadOrNew()

					Expect(err).To(BeNil())
					Expect(cfg.ClientMods).To(HaveLen(len(TestingConfig.ClientMods)))

					for i := 0; i < len(cfg.ClientMods); i++ {
						cfgMod := cfg.ClientMods[i]
						xpctMod := TestingConfig.ClientMods[i]
						Expect(cfgMod.CliName).To(Equal(xpctMod.CliName))
						Expect(cfgMod.FriendlyName).To(Equal(xpctMod.FriendlyName))
						Expect(cfgMod.Description).To(Equal(xpctMod.Description))
						Expect(cfgMod.DetailsURL).To(Equal(xpctMod.DetailsURL))
						Expect(cfgMod.LatestURL).To(Equal(xpctMod.LatestURL))
					}

					Expect(cfg.ModInstallations).To(HaveLen(len(TestingConfig.ModInstallations)))
					for cliName, cfgMod := range cfg.ModInstallations {
						xpctMod := TestingConfig.ModInstallations[cliName]
						Expect(cfgMod.DownloadURL).To(Equal(xpctMod.DownloadURL))
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
				afero.WriteFile(fs, configPath, []byte("{"), 0644)

				cfg, err := configIo.LoadOrNew()

				Expect(cfg).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("Save func", func() {
		It("saves empty config", func() {
			err := configIo.Save(&mc.UserModConfig{})

			Expect(err).To(BeNil())

			b, err := afero.ReadFile(fs, configPath)

			Expect(string(b)).To(Equal(emptyContent))
		})

		It("saves non-empty config", func() {
			err := configIo.Save(TestingConfig)

			Expect(err).To(BeNil())

			b, err := afero.ReadFile(fs, configPath)

			Expect(string(b)).To(Equal(content))
		})
	})

	Context("New func", func() {
		It("returns an empty config", func() {
			cfg := mc.NewUserModConfig()

			Expect(cfg).ToNot(BeNil())
			Expect(cfg.ClientMods).ToNot(BeNil())
			Expect(cfg.ClientMods).To(BeEmpty())
			Expect(cfg.ModInstallations).ToNot(BeNil())
			Expect(cfg.ModInstallations).To(BeEmpty())
		})
	})
})

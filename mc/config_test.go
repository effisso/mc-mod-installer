package mc_test

import (
	"encoding/json"
	"mcmods/mc"
	. "mcmods/testdata"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Config IO", func() {
	var fs afero.Fs

	mcPath := "/path/to/mc"

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		fs = afero.NewMemMapFs()
		mc.ViperInstance.Set(mc.InstallPathKey, mcPath)
	})

	Context("Client config", func() {
		var configIo mc.ModConfigIo
		var configPath, content, emptyContent string

		BeforeEach(func() {
			configPath = filepath.Join(mcPath, mc.ModFolderName, mc.ModConfigFileName)

			b, _ := json.MarshalIndent(TestingConfig, "", "\t")
			content = string(b)

			b, _ = json.MarshalIndent(mc.UserModConfig{}, "", "\t")
			emptyContent = string(b)

			mcfs := &mc.LocalFileSystem{Fs: fs}
			configIo = mc.NewUserModConfigIo(mcfs)
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

	Context("Server config saver", func() {
		var saver mc.ServerConfigSaver

		wd, wdErr := os.Getwd()
		Expect(wdErr).To(BeNil())
		cfgPath := filepath.Join(wd, "mc", "server_mods.json")

		BeforeEach(func() {
			saver = mc.NewServerConfigSaver(fs)
		})

		It("fails if the file doesn't already exist", func() {
			err := saver.Save()

			Expect(err).ToNot(BeNil())
		})

		It("overwrites the file if it exists", func() {
			Expect(afero.WriteFile(fs, cfgPath, []byte{}, 0644)).To(BeNil())

			err := saver.Save()

			Expect(err).To(BeNil())
			expected, err := json.MarshalIndent(TestingServerGroups, "", "\t")
			Expect(err).To(BeNil())
			actual, err := afero.ReadFile(fs, cfgPath)
			Expect(err).To(BeNil())
			Expect(string(actual)).To(Equal(string(expected)))
		})
	})
})

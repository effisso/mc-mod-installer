package mc_test

import (
	"errors"
	"fmt"
	"mcmods/mc"
	. "mcmods/testdata"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Installer", func() {
	var installer mc.ModInstaller
	var cfg *mc.ClientModConfig
	var singleMod []*mc.Mod

	installLoc := "/test/path"

	BeforeEach(func() {
		InitTestData()
		singleMod = []*mc.Mod{TestingClientMod1}
		mc.ServerGroups = TestingServerGroups
		installer = mc.NewModInstaller()
		cfg = &mc.ClientModConfig{
			ModInstallations: map[string]mc.ModInstallation{},
			ClientMods:       []*mc.Mod{},
		}

		mc.ViperInstance.Set(mc.InstallPathKey, installLoc)
	})

	It("passes correct args to the downloader", func() {
		dl := verifyingDownloader{
			ExpectedPath: filepath.Join(installLoc, mc.ModsFolderName, TestingClientMod1.CliName+".jar"),
			ExpectedMod:  TestingClientMod1,
		}

		err := installer.InstallMods(dl, singleMod, cfg)

		Expect(err).To(BeNil())
	})

	It("calls the downloader the correct number of times", func() {
		mods := []*mc.Mod{TestingClientMod1, TestingClientMod2, TestingServerRequired1}
		dl := &countingDownloader{}

		err := installer.InstallMods(dl, mods, cfg)

		Expect(err).To(BeNil())
		Expect(dl.CallCount).To(Equal(len(mods)))
	})

	It("returns errors thrown by the downloader", func() {
		dl := emptyDownloader{Err: errors.New("test")}

		err := installer.InstallMods(dl, singleMod, cfg)

		Expect(err).To(Equal(dl.Err))
	})

	It("adds install items to the config", func() {
		mods := []*mc.Mod{TestingClientMod1, TestingClientMod2}
		dl := &emptyDownloader{}
		nowText := fmt.Sprint(time.Now().Format(time.UnixDate))

		err := installer.InstallMods(dl, mods, cfg)

		Expect(err).To(BeNil())
		Expect(cfg.ModInstallations).To(HaveLen(len(mods)))

		verifyInstall(TestingClientMod1, cfg, nowText)
		verifyInstall(TestingClientMod2, cfg, nowText)
	})
})

func verifyInstall(mod *mc.Mod, cfg *mc.ClientModConfig, nowText string) {
	install := cfg.ModInstallations[mod.CliName]
	Expect(install.DownloadUrl).To(Equal(mod.LatestUrl))
	Expect(install.Timestamp).To(Equal(nowText))
}

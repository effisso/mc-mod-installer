package cmd_test

import (
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Visit Cmd", func() {
	var fs afero.Fs
	var launcher *spyLauncher
	var mapValidator *nameMapperValidator
	var cfgIoSpy *clientConfigIoSpy

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		mapb := false
		mapValidator = &nameMapperValidator{
			ClientMods: TestingClientMods,
			Visited:    &mapb,
			fakeNameMapper: fakeNameMapper{
				Map: TestingCliModMap,
			},
		}
		cmd.NameMapper = mapValidator

		cfgIoSpy = &clientConfigIoSpy{
			LoadReturn: TestingConfig,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}

		launcher = &spyLauncher{}
		cmd.BrowserLauncher = launcher
	})

	It("launches a browser with the mod's details url", func() {
		launcher.expectedURL = TestingClientMod1.DetailsURL
		cmd.RootCmd.SetArgs([]string{"visit", TestingClientMod1.CliName})

		err := cmd.RootCmd.Execute()

		Expect(err).To(BeNil())
		Expect(launcher.launched).To(BeTrue())
	})

	It("returns an error for an unknown mod", func() {
		cmd.RootCmd.SetArgs([]string{"visit", "unknown"})

		err := cmd.RootCmd.Execute()

		Expect(err).ToNot(BeNil())
		Expect(launcher.launched).To(BeFalse())
	})
})

type spyLauncher struct {
	launched    bool
	expectedURL string
	Return      error
}

func (l *spyLauncher) Open(url string) error {
	l.launched = true
	Expect(url).To(Equal(l.expectedURL))
	return l.Return
}

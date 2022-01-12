package cmd_test

import (
	"mcmods/cmd"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Visit Cmd", func() {
	var mapValidator *nameMapperValidator
	var launcher *spyLauncher

	BeforeEach(func() {
		rootCmdTestSetup()

		mapb := false
		mapValidator = &nameMapperValidator{
			ClientMods: TestingClientMods,
			Visited:    &mapb,
			fakeNameMapper: fakeNameMapper{
				Map: TestingCliModMap,
			},
		}
		cmd.NameMapper = mapValidator

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

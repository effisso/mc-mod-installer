package cmd_test

import (
	"errors"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestCmds(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var _ = Describe("Install Cmd", func() {
	var fs afero.Fs
	var cfgIoFake *clientConfigIoSpy
	var dl mc.ModDownloader

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		cmd.ResetInstallVars()
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		b := false
		cfgIoFake = &clientConfigIoSpy{
			Saved:      &b,
			LoadReturn: TestingConfig,
		}

		dl = &fakeDownloader{}
		cmd.CreateDownloaderFunc = func(fs mc.FileSystem) mc.ModDownloader {
			return dl
		}

		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoFake
		}
	})

	Context("verify filter, install, and save config", func() {
		var verifyFilter *filterVerifier
		var verifyInstaller *installerVerifier

		BeforeEach(func() {
			af := false
			// verifies args passed to the filter
			verifyFilter = &filterVerifier{
				XGroups: []string{},
				XMods:   []string{},
				Cfg:     TestingConfig,
				Force:   false,
				Visited: &af,
				emptyFilter: emptyFilter{
					Return: []*mc.Mod{},
				},
			}

			bf := false
			// verifies args passed to the installer
			verifyInstaller = &installerVerifier{
				Downloader: dl,
				Cfg:        TestingConfig,
				Mods:       []*mc.Mod{},
				Visited:    &bf,
				emptyInstaller: emptyInstaller{
					Return: nil,
				},
			}

			cmd.Filter = verifyFilter
			cmd.Installer = verifyInstaller
		})

		When("no args", func() {
			It("adds server-only to the list of groups to exclude", func() {
				verifyFilter.XGroups = []string{"server-only"}
				verifyFilter.Return = append(TestingClientMods, TestingServerOptional1, TestingServerPerformance1, TestingServerRequired1)
				verifyInstaller.Mods = verifyFilter.Return
				cmd.RootCmd.SetArgs([]string{"install"})

				err := cmd.RootCmd.Execute()

				Expect(err).To(BeNil(), "no error should have been returned")
				Expect(*verifyFilter.Visited).To(BeTrue(), "mods not filtered")
				Expect(*verifyInstaller.Visited).To(BeTrue(), "mods not installed")
				Expect(*cfgIoFake.Saved).To(BeTrue())
			})
		})

		When("--full-server", func() {
			It("installs all server groups", func() {
				cmd.RootCmd.SetArgs([]string{"install", "--full-server"})

				err := cmd.RootCmd.Execute()

				Expect(err).To(BeNil(), "no error should have been returned")
				Expect(*verifyFilter.Visited).To(BeTrue(), "mods not filtered")
				Expect(*verifyInstaller.Visited).To(BeTrue(), "mods not installed")
				Expect(*cfgIoFake.Saved).To(BeTrue())
			})
		})

		When("--client-only", func() {
			It("adds all server groups to the exclude list", func() {
				verifyFilter.XGroups = TestingServerGroupNames
				cmd.RootCmd.SetArgs([]string{"install", "--client-only"})

				err := cmd.RootCmd.Execute()

				Expect(err).To(BeNil(), "no error should have been returned")
				Expect(*verifyFilter.Visited).To(BeTrue(), "mods not filtered")
				Expect(*verifyInstaller.Visited).To(BeTrue(), "mods not installed")
				Expect(*cfgIoFake.Saved).To(BeTrue())
			})
		})

		It("returns error from save", func() {
			cfgIoFake.SaveErr = errors.New("save err")
			cmd.RootCmd.SetArgs([]string{"install", "--full-server"})

			err := cmd.RootCmd.Execute()
			Expect(err).To(Equal(cfgIoFake.SaveErr))
		})
	})
	Context("client install", func() {
		var filter *filterVerifier
		var emInstaller *emptyInstaller

		BeforeEach(func() {
			filter = &filterVerifier{}
			emInstaller = &emptyInstaller{}

			af := false
			// verifies args passed to the filter
			filter = &filterVerifier{
				XGroups: TestingServerGroupNames,
				XMods:   []string{},
				Cfg:     TestingConfig,
				Force:   false,
				Visited: &af,
				emptyFilter: emptyFilter{
					Return: TestingClientMods,
				},
			}

			cmd.Filter = filter
			cmd.Installer = emInstaller
		})

		Context("client only", func() {
			When("true", func() {
				BeforeEach(func() {
					af := false
					// verifies args passed to the filter
					filter = &filterVerifier{
						XGroups: TestingServerGroupNames,
						XMods:   []string{},
						Cfg:     TestingConfig,
						Force:   false,
						Visited: &af,
						emptyFilter: emptyFilter{
							Return: TestingClientMods,
						},
					}

					cmd.Filter = filter
				})

				It("filters out all server groups", func() {
					cmd.RootCmd.SetArgs([]string{"install", "--client-only"})

					err := cmd.RootCmd.Execute()

					Expect(err).To(BeNil(), "no error should have been returned")
					Expect(*filter.Visited).To(BeTrue(), "mods not filtered")
				})
			})
			When("false", func() {
				BeforeEach(func() {
					af := false
					// verifies args passed to the filter
					filter = &filterVerifier{
						XGroups: TestingServerGroupNames,
						XMods:   []string{},
						Cfg:     TestingConfig,
						Force:   false,
						Visited: &af,
						emptyFilter: emptyFilter{
							Return: TestingClientMods,
						},
					}

					cmd.Filter = filter
				})

				It("adds server-only to the group exclusion list", func() {
					filter.XGroups = []string{"performance", "server-only"}
					cmd.RootCmd.SetArgs([]string{"install", "--x-group", "performance"})

					err := cmd.RootCmd.Execute()

					Expect(err).To(BeNil(), "no error should have been returned")
					Expect(*filter.Visited).To(BeTrue(), "mods not filtered")
				})

				It("excludes the mods exclusion list", func() {
					filter.XGroups = []string{"server-only"}
					filter.XMods = []string{TestingClientMod1.CliName}
					cmd.RootCmd.SetArgs([]string{"install", "--x-mod", TestingClientMod1.CliName})

					err := cmd.RootCmd.Execute()

					Expect(err).To(BeNil(), "no error should have been returned")
					Expect(*filter.Visited).To(BeTrue(), "mods not filtered")
				})
			})
		})
	})
})

// ----
// Mock Filters
// ----

// just return some mods
type emptyFilter struct {
	Return []*mc.Mod
	Err    error
}

func (f emptyFilter) FilterAllMods(xGroups []string, xMods []string, cfg *mc.UserModConfig, force bool) ([]*mc.Mod, error) {
	return f.Return, f.Err
}

// verify the filter arguments
type filterVerifier struct {
	emptyFilter
	XGroups []string
	XMods   []string
	Cfg     *mc.UserModConfig
	Force   bool
	Visited *bool
}

func (f filterVerifier) FilterAllMods(xGroups []string, xMods []string, cfg *mc.UserModConfig, force bool) ([]*mc.Mod, error) {
	*(f.Visited) = true
	Expect(xGroups).To(ConsistOf(f.XGroups))
	Expect(xMods).To(ConsistOf(f.XMods))
	Expect(cfg).To(Equal(f.Cfg))
	Expect(force).To(Equal(f.Force))
	return f.Return, f.Err
}

// ----
// Mock Installers
// ----

// just return
type emptyInstaller struct {
	Return error
}

func (i emptyInstaller) InstallMods(downloader mc.ModDownloader, mods []*mc.Mod, cfg *mc.UserModConfig) error {
	return i.Return
}

// verify arguments
type installerVerifier struct {
	emptyInstaller
	Downloader mc.ModDownloader
	Mods       []*mc.Mod
	Cfg        *mc.UserModConfig
	Visited    *bool
}

func (i installerVerifier) InstallMods(downloader mc.ModDownloader, mods []*mc.Mod, cfg *mc.UserModConfig) error {
	*(i.Visited) = true
	Expect(downloader).To(Equal(i.Downloader))
	Expect(mods).To(ConsistOf(i.Mods))
	Expect(cfg).To(Equal(i.Cfg))
	return i.Return
}

// ----
// ConfigIo
// ----

type clientConfigIoSpy struct {
	LoadReturn *mc.UserModConfig
	LoadErr    error
	Saved      *bool
	SaveErr    error
}

func (i clientConfigIoSpy) LoadOrNew() (*mc.UserModConfig, error) {
	return i.LoadReturn, i.LoadErr
}

func (i clientConfigIoSpy) Save(cfg *mc.UserModConfig) error {
	*(i.Saved) = true
	return i.SaveErr
}

// ----
// Downloader
// ----

type fakeDownloader struct{}

func (fakeDownloader) Download(mod *mc.Mod, relPath string) error {
	return nil
}

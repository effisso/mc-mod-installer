package cmd_test

import (
	"bytes"
	"errors"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("List Cmd", func() {
	var fs afero.Fs
	var outBuffer *bytes.Buffer
	var clientModOutput, serverModOutput string
	var nameValidator *vNameValidator
	var cfgIoFake *clientConfigIoSpy

	BeforeEach(func() {
		InitTestData()
		clientModOutput = strings.Join(TestingClientCliNames, "\n") + "\n"
		serverModOutput = strings.Join(TestingServerCliNames, "\n") + "\n"

		cmd.ResetListVars()

		mc.ServerGroups = TestingServerGroups
		cmd.ResetInstallVars()
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		cfg := TestingConfig
		cfg.ClientMods = TestingClientMods
		cfgIoFake = &clientConfigIoSpy{
			LoadReturn: cfg,
		}

		groupb := false
		modb := false
		nameValidator = &vNameValidator{
			Groups:             []string{},
			Mods:               []string{},
			VisitedGroup:       &groupb,
			VisitedMod:         &modb,
			Map:                TestingCliModMap,
			emptyNameValidator: emptyNameValidator{},
		}

		cmd.NameValidator = nameValidator
		cmd.ConfigIo = cfgIoFake

		outBuffer = bytes.NewBufferString("")

		cmd.RootCmd.SetOut(outBuffer)
	})

	Context("mods", func() {
		Context("client only", func() {
			It("prints all client mods, and no others", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--client"})

				executeAndVerifyOutput(outBuffer, clientModOutput, false)
			})
		})

		Context("server only", func() {
			It("prints all server mods, and no others", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--server"})

				executeAndVerifyOutput(outBuffer, serverModOutput, false)
			})
		})

		Context("client and server", func() {
			When("explicitly provided", func() {
				It("prints both client and server mods", func() {
					cmd.RootCmd.SetArgs([]string{"list", "mods", "--client", "--server"})

					executeAndVerifyOutput(outBuffer, serverModOutput+clientModOutput, false)
				})
			})

			When("both undeclared", func() {
				It("prints both client and server mods", func() {
					cmd.RootCmd.SetArgs([]string{"list", "mods"})

					executeAndVerifyOutput(outBuffer, serverModOutput+clientModOutput, false)
				})
			})
		})

		Context("installed", func() {
			It("shows only installed mods", func() {
				installedNameSlice := make([]string, 0, 2)
				for name := range TestingConfig.ModInstallations {
					installedNameSlice = append(installedNameSlice, name)
				}
				installedNamesOutput := strings.Join(installedNameSlice, "\n") + "\n"
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed"})

				executeAndVerifyOutput(outBuffer, installedNamesOutput, false)
			})

			When("combined with target switch", func() {
				Context("client", func() {
					It("displays only installed client mods", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--client"})

						executeAndVerifyOutput(outBuffer, TestingClientMod1.CliName+"\n", false)
					})
				})
				Context("server", func() {
					It("displays only installed server mods", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--server"})

						executeAndVerifyOutput(outBuffer, TestingServerRequired1.CliName+"\n", false)
					})
				})
			})
		})

		Context("not installed", func() {
			It("shows only mods which are not installed", func() {
				notInstalledNameSlice := make([]string, 0, 5)
				for name := range TestingCliModMap {
					if _, installed := TestingConfig.ModInstallations[name]; !installed {
						notInstalledNameSlice = append(notInstalledNameSlice, name)
					}
				}
				notInstalledNamesOutput := strings.Join(notInstalledNameSlice, "\n") + "\n"
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--not-installed"})

				executeAndVerifyOutput(outBuffer, notInstalledNamesOutput, false)
			})

			When("combined with target switch", func() {
				Context("client", func() {
					It("displays only client mods which are not installed", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--not-installed", "--client"})

						executeAndVerifyOutput(outBuffer, TestingClientMod2.CliName+"\n", false)
					})
				})
				Context("server", func() {
					It("displays only server mods which are not installed", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--not-installed", "--server"})

						expectedOutput := strings.Join([]string{
							TestingServerOnly1.CliName,
							TestingServerOptional1.CliName,
							TestingServerPerformance1.CliName,
						}, "\n") + "\n"

						executeAndVerifyOutput(outBuffer, expectedOutput, false)
					})
				})
			})

			It("prints all client and server mods if both install switches provided", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--not-installed"})

				executeAndVerifyOutput(outBuffer, serverModOutput+clientModOutput, false)
			})
		})

		Context("group", func() {
			requiredGroupName := "required"
			invalidGroup := "invalid"

			It("returns an error when the group name is invalid", func() {
				nameValidator.Groups = []string{invalidGroup}
				nameValidator.GroupsReturn = errors.New("invalid group name")
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", invalidGroup})

				err := cmd.RootCmd.Execute()

				Expect(err).To(Equal(nameValidator.GroupsReturn))
			})

			It("shows only mods from the specified server group", func() {
				nameValidator.Groups = []string{requiredGroupName}
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", requiredGroupName})

				executeAndVerifyOutput(outBuffer, TestingServerRequired1.CliName+"\n", false)
				Expect(*nameValidator.VisitedGroup).To(BeTrue())
			})

			It("returns an error when combined with the client switch", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", "doesnt-matter", "--client"})

				err := cmd.RootCmd.Execute()

				Expect(err).ToNot(BeNil())
			})

			It("displays only mods in the group, even with the server switch", func() {
				nameValidator.Groups = []string{requiredGroupName}
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", requiredGroupName, "--server"})

				executeAndVerifyOutput(outBuffer, TestingServerRequired1.CliName+"\n", false)
			})
		})
	})

	Context("groups", func() {
		It("prints all server groups", func() {
			cmd.RootCmd.SetArgs([]string{"list", "groups"})

			executeAndVerifyOutput(outBuffer, strings.Join(TestingServerGroupNames, "\n")+"\n", false)
		})
	})
})

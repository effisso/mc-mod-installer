package cmd_test

import (
	"mcmods/cmd"
	. "mcmods/testdata"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List Cmd", func() {
	var td *rootTestData
	var clientModOutput, serverModOutput string

	BeforeEach(func() {
		td = rootCmdTestSetup()

		clientModOutput = strings.Join(TestingClientCliNames, "\n") + "\n"
		serverModOutput = strings.Join(TestingServerCliNames, "\n") + "\n"
	})

	Context("mods", func() {
		Context("client only", func() {
			It("prints all client mods, and no others", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--client"})

				executeAndVerifyOutput(td.outBuffer, clientModOutput, false)
			})
		})

		Context("server only", func() {
			It("prints all server mods, and no others", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--server"})

				executeAndVerifyOutput(td.outBuffer, serverModOutput, false)
			})
		})

		Context("client and server", func() {
			When("explicitly provided", func() {
				It("prints both client and server mods", func() {
					cmd.RootCmd.SetArgs([]string{"list", "mods", "--client", "--server"})

					executeAndVerifyOutput(td.outBuffer, serverModOutput+clientModOutput, false)
				})
			})

			When("both undeclared", func() {
				It("prints both client and server mods", func() {
					cmd.RootCmd.SetArgs([]string{"list", "mods"})

					executeAndVerifyOutput(td.outBuffer, serverModOutput+clientModOutput, false)
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

				executeAndVerifyOutput(td.outBuffer, installedNamesOutput, false)
			})

			When("combined with target switch", func() {
				Context("client", func() {
					It("displays only installed client mods", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--client"})

						executeAndVerifyOutput(td.outBuffer, TestingClientMod1.CliName+"\n", false)
					})
				})
				Context("server", func() {
					It("displays only installed server mods", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--server"})

						executeAndVerifyOutput(td.outBuffer, TestingServerRequired1.CliName+"\n", false)
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

				executeAndVerifyOutput(td.outBuffer, notInstalledNamesOutput, false)
			})

			When("combined with target switch", func() {
				Context("client", func() {
					It("displays only client mods which are not installed", func() {
						cmd.RootCmd.SetArgs([]string{"list", "mods", "--not-installed", "--client"})

						executeAndVerifyOutput(td.outBuffer, TestingClientMod2.CliName+"\n", false)
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

						executeAndVerifyOutput(td.outBuffer, expectedOutput, false)
					})
				})
			})

			It("prints all client and server mods if both install switches provided", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--installed", "--not-installed"})

				executeAndVerifyOutput(td.outBuffer, serverModOutput+clientModOutput, false)
			})
		})

		Context("group", func() {
			requiredGroupName := "required"
			invalidGroup := "invalid"

			It("returns an error when the group name is invalid", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", invalidGroup})

				err := cmd.RootCmd.Execute()

				Expect(err).ToNot(BeNil())
			})

			It("shows only mods from the specified server group", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", requiredGroupName})

				executeAndVerifyOutput(td.outBuffer, TestingServerRequired1.CliName+"\n", false)
			})

			It("returns an error when combined with the client switch", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", "doesnt-matter", "--client"})

				err := cmd.RootCmd.Execute()

				Expect(err).ToNot(BeNil())
			})

			It("displays only mods in the group, even with the server switch", func() {
				cmd.RootCmd.SetArgs([]string{"list", "mods", "--group", requiredGroupName, "--server"})

				executeAndVerifyOutput(td.outBuffer, TestingServerRequired1.CliName+"\n", false)
			})
		})
	})

	Context("groups", func() {
		It("prints all server groups", func() {
			cmd.RootCmd.SetArgs([]string{"list", "groups"})

			executeAndVerifyOutput(td.outBuffer, strings.Join(TestingServerGroupNames, "\n")+"\n", false)
		})
	})
})

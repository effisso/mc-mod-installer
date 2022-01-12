package cmd_test

import (
	"mcmods/cmd"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Version Cmd", func() {
	var td *rootTestData

	curVer := "1.2.3"

	BeforeEach(func() {
		td = rootCmdTestSetup()

		cmd.RootCmd.SetArgs([]string{"version"})
		cmd.ToolVersion = curVer
	})

	It("prints the current version to the user", func() {
		executeAndVerifyOutput(td.outBuffer, curVer, true)
	})
})

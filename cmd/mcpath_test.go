package cmd_test

import (
	"mcmods/cmd"
	"mcmods/mc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MC Path Cmd", func() {
	var td *rootTestData

	mcPathVal := "/some/path/to/.minecraft"

	BeforeEach(func() {
		td = rootCmdTestSetup()

		cmd.ViperInstance.Set(mc.InstallPathKey, mcPathVal)
	})

	It("no arguments - prints the path in viper to the user", func() {
		cmd.RootCmd.SetArgs([]string{"mcpath"})

		executeAndVerifyOutput(td.outBuffer, mcPathVal, true)
	})

	It("set - updates viper with the new value", func() {
		updatedPath := "/a/different/path/to/.minecraft"
		cmd.RootCmd.SetArgs([]string{"mcpath", "--set", updatedPath})

		cmd.RootCmd.Execute()

		pathInViper := cmd.ViperInstance.GetString(mc.InstallPathKey)
		Expect(pathInViper).To(Equal(updatedPath))
	})
})

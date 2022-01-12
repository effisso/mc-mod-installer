package cmd_test

import (
	"mcmods/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Docs Cmd", func() {
	var launcher *spyLauncher

	BeforeEach(func() {
		rootCmdTestSetup()

		launcher = &spyLauncher{}
		cmd.BrowserLauncher = launcher
	})

	It("launches a browser with the doc url", func() {
		launcher.expectedURL = cmd.DocURL
		cmd.RootCmd.SetArgs([]string{"docs"})

		err := cmd.RootCmd.Execute()

		Expect(err).To(BeNil())
		Expect(launcher.launched).To(BeTrue())
	})
})

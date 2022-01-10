package cmd_test

import (
	"bytes"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	"github.com/spf13/afero"
)

var _ = Describe("Version Cmd", func() {
	var fs afero.Fs
	var cfgIoSpy *clientConfigIoSpy
	var outBuffer *bytes.Buffer

	curVer := "1.2.3"

	BeforeEach(func() {
		mc.ServerGroups = TestingServerGroups
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		cfgIoSpy = &clientConfigIoSpy{
			LoadReturn: TestingConfig,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}

		outBuffer = bytes.NewBufferString("")

		cmd.RootCmd.SetOut(outBuffer)
		cmd.RootCmd.SetArgs([]string{"version"})
		cmd.ToolVersion = curVer
	})

	It("prints the current version to the user", func() {
		executeAndVerifyOutput(outBuffer, curVer, true)
	})
})

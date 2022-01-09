package cmd_test

import (
	"bytes"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("MC Path Cmd", func() {
	var fs afero.Fs
	var cfgIoSpy *clientConfigIoSpy
	var outBuffer *bytes.Buffer

	mcPathVal := "/some/path/to/.minecraft"

	BeforeEach(func() {
		cmd.ResetPathVars()
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		cmd.ViperInstance.Set(mc.InstallPathKey, mcPathVal)

		cfgIoSpy = &clientConfigIoSpy{
			LoadReturn: TestingConfig,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}

		outBuffer = bytes.NewBufferString("")

		cmd.RootCmd.SetOut(outBuffer)
	})

	It("no arguments - prints the path in viper to the user", func() {
		cmd.RootCmd.SetArgs([]string{"mcpath"})

		executeAndVerifyOutput(outBuffer, mcPathVal+"\n", true)
	})

	It("set - updates viper with the new value", func() {
		updatedPath := "/a/different/path/to/.minecraft"
		cmd.RootCmd.SetArgs([]string{"mcpath", "--set", updatedPath})

		cmd.RootCmd.Execute()

		pathInViper := cmd.ViperInstance.GetString(mc.InstallPathKey)
		Expect(pathInViper).To(Equal(updatedPath))
	})
})

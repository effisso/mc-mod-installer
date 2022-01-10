package cmd_test

import (
	"bytes"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var _ = Describe("Describe Cmd", func() {
	var fs afero.Fs
	var outBuffer *bytes.Buffer

	var cfgIoSpy *clientConfigIoSpy

	BeforeEach(func() {
		cmd.ResetVars()

		fs = afero.NewMemMapFs()
		mc.ServerGroups = TestingServerGroups

		cmd.CreateFsFunc = func(ftpArgs *mc.FtpArgs) (mc.FileSystem, error) {
			return mc.LocalFileSystem{Fs: fs}, nil
		}

		cfgIoSpy = &clientConfigIoSpy{
			LoadReturn: TestingConfig,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}

		outBuffer = bytes.NewBufferString("")

		cmd.RootCmd.SetOut(outBuffer)
	})

	Context("FTP", func() {
		expectedUsername := "hackerman"
		expectedServer := "all-your-base-are-belong-to-us"
		expectedPassword := "happyhippo12"

		BeforeEach(func() {
			cmd.RootCmd.SetArgs([]string{"--user", expectedUsername, "--password", expectedPassword, "--ftp-server", expectedServer})

			// set Run to an empty function just so the persistent pre/post-run funcs also execute
			cmd.RootCmd.Run = func(cmd *cobra.Command, args []string) {}
		})

		AfterEach(func() {
			cmd.RootCmd.Run = nil
		})

		It("updates viper with ftp user and server args", func() {
			err := cmd.RootCmd.Execute()

			Expect(err).To(BeNil())
			Expect(cmd.ViperInstance.GetString(mc.FtpServerKey)).To(Equal(expectedServer))
			Expect(cmd.ViperInstance.GetString(mc.FtpUserKey)).To(Equal(expectedUsername))
		})

		It("passes ftp args to the CreateFsFunc", func() {
			cmd.CreateFsFunc = func(ftpArgs *mc.FtpArgs) (mc.FileSystem, error) {
				Expect(ftpArgs.User).To(Equal(expectedUsername))
				Expect(ftpArgs.Pw).To(Equal(expectedPassword))
				Expect(ftpArgs.Server).To(Equal(expectedServer))
				return mc.LocalFileSystem{Fs: fs}, nil
			}

			err := cmd.RootCmd.Execute()

			Expect(err).To(BeNil())
		})
	})
})

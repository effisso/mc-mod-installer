package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mcmods/cmd"
	"mcmods/mc"
	. "mcmods/testdata"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Describe Cmd", func() {
	var fs afero.Fs
	var outBuffer *bytes.Buffer

	var mapValidator *nameMapperValidator
	var cfgIoSpy *clientConfigIoSpy

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		cmd.ResetInstallVars()
		fs = afero.NewMemMapFs()
		cmd.ViperInstance.SetFs(fs)

		mapb := false
		mapValidator = &nameMapperValidator{
			ClientMods: TestingClientMods,
			Visited:    &mapb,
			fakeNameMapper: fakeNameMapper{
				Map: TestingCliModMap,
			},
		}
		cmd.NameMapper = mapValidator

		cfgIoSpy = &clientConfigIoSpy{
			LoadReturn: TestingConfig,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}

		outBuffer = bytes.NewBufferString("")

		cmd.RootCmd.SetOut(outBuffer)
	})

	Context("mod", func() {
		It("no mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "mod"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("invalid mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "mod", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("describes the mod", func() {
			m := TestingClientMod1
			expectedOutput := fmt.Sprintf("\n%s (%s)\n-----\n%s\nWebsite:  %s\nLatest package:  %s",
				m.FriendlyName, m.CliName, m.Description, m.DetailsURL, m.LatestURL)

			cmd.RootCmd.SetArgs([]string{"describe", "mod", m.CliName})

			executeAndVerifyOutput(outBuffer, expectedOutput, true)
		})
	})

	Context("group", func() {
		It("no mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "group"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("invalid group name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "group", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("describes the group", func() {
			m1 := TestingServerRequired1
			m2 := &mc.Mod{
				CliName:     "second",
				Description: "used to verify that multiple groups can be printed",
			}
			TestingServerGroups["required"].Mods = append(TestingServerGroups["required"].Mods, m2)

			cmd.RootCmd.SetArgs([]string{"describe", "group", "required"})

			executeAndVerifyOutput(outBuffer, m1.CliName+"\n"+m2.CliName, false)
		})
	})

	Context("install", func() {
		It("no mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "install"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("invalid mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "install", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("describes the install", func() {
			m := TestingClientMod1
			expectedOutput := fmt.Sprintf("\n%s (%s)\n-----\nInstall timestamp:  %s\nUp-to-date:  %t",
				m.FriendlyName, m.CliName, "123", false)

			cmd.RootCmd.SetArgs([]string{"describe", "install", TestingClientMod1.CliName})

			executeAndVerifyOutput(outBuffer, expectedOutput, true)
		})

		It("informs when not installed", func() {
			expectedOutput := fmt.Sprintf("Not Installed.")

			cmd.RootCmd.SetArgs([]string{"describe", "install", TestingServerOnly1.CliName})

			executeAndVerifyOutput(outBuffer, expectedOutput, true)
		})
	})

	Context("other", func() {
		It("returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "invalid", "doesnt-matter"})

			Expect(cmd.RootCmd.Execute()).To(Not(BeNil()))
		})
	})
})

func executeAndVerifyOutput(outBuffer io.Reader, expectedOutput string, lineOrderMatters bool) {
	err := cmd.RootCmd.Execute()

	Expect(err).To(BeNil())

	out, err := ioutil.ReadAll(outBuffer)

	Expect(err).To(BeNil())
	strOut := string(out)

	if lineOrderMatters {
		Expect(strOut).To(Equal(expectedOutput))
	} else {
		newlineSplitFunc := func(c rune) bool {
			return c == '\n'
		}
		outLines := strings.FieldsFunc(strOut, newlineSplitFunc)
		expectedLines := strings.FieldsFunc(expectedOutput, newlineSplitFunc)

		Expect(outLines).To(ConsistOf(expectedLines))
	}
}

// ----
// Name Mapper Mocks
// ----

type fakeNameMapper struct {
	Map mc.ModMap
}

func (m fakeNameMapper) MapAllMods(clientMods []*mc.Mod) mc.ModMap {
	return m.Map
}

type nameMapperValidator struct {
	fakeNameMapper
	ClientMods []*mc.Mod
	Visited    *bool
}

func (m nameMapperValidator) MapAllMods(clientMods []*mc.Mod) mc.ModMap {
	*(m.Visited) = true
	Expect(clientMods).To(ConsistOf(m.ClientMods))
	return m.Map
}

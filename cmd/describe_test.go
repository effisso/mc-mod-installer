package cmd_test

import (
	"bytes"
	"errors"
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
	var nameValidator *vNameValidator
	var cfgIoFake *clientConfigIoSpy

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

		cfg := TestingConfig
		cfg.ClientMods = TestingClientMods
		cfgIoFake = &clientConfigIoSpy{
			LoadReturn: cfg,
		}
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoFake
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
			nameValidator.Mods = []string{"invalid"}
			nameValidator.ModsReturn = errors.New("invalid mod name")
			cmd.RootCmd.SetArgs([]string{"describe", "mod", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(nameValidator.ModsReturn))
		})

		It("describes the mod", func() {
			m := TestingClientMod1
			nameValidator.Mods = []string{m.CliName}
			expectedOutput := fmt.Sprintf("\n%s (%s)\n-----\n%s\nWebsite:  %s\nLatest package:  %s\n",
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
			nameValidator.Groups = []string{"invalid"}
			nameValidator.GroupsReturn = errors.New("invalid group name")
			cmd.RootCmd.SetArgs([]string{"describe", "group", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(nameValidator.GroupsReturn))
		})

		It("describes the group", func() {
			m := TestingServerRequired1
			nameValidator.Groups = []string{"required"}
			expectedOutput := fmt.Sprintf("%s\n", m.CliName)

			cmd.RootCmd.SetArgs([]string{"describe", "group", "required"})

			executeAndVerifyOutput(outBuffer, expectedOutput, false)
		})
	})

	Context("install", func() {
		It("no mod name returns an error", func() {
			cmd.RootCmd.SetArgs([]string{"describe", "install"})

			err := cmd.RootCmd.Execute()

			Expect(err).ToNot(BeNil())
		})

		It("invalid mod name returns an error", func() {
			nameValidator.Mods = []string{"invalid"}
			nameValidator.ModsReturn = errors.New("invalid mod name")
			cmd.RootCmd.SetArgs([]string{"describe", "install", "invalid"})

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(nameValidator.ModsReturn))
		})

		It("describes the install", func() {
			m := TestingClientMod1
			nameValidator.Mods = []string{TestingClientMod1.CliName}
			expectedOutput := fmt.Sprintf("\n%s (%s)\n-----\nInstall timestamp:  %s\nUp-to-date:  %t\n",
				m.FriendlyName, m.CliName, "123", false)

			cmd.RootCmd.SetArgs([]string{"describe", "install", TestingClientMod1.CliName})

			executeAndVerifyOutput(outBuffer, expectedOutput, true)
		})

		It("informs when not installed", func() {
			nameValidator.Mods = []string{TestingServerOnly1.CliName}
			expectedOutput := fmt.Sprintf("Not Installed.\n")

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
		outLines := strings.Split(strOut, "\n")
		expectedLines := strings.Split(expectedOutput, "\n")

		Expect(outLines).To(ConsistOf(expectedLines))
	}
}

// ----
// Name Validator Mocks
// ----

type emptyNameValidator struct {
	GroupsReturn error
	ModsReturn   error
}

func (v emptyNameValidator) ValidateServerGroups(groups []string) error {
	return v.GroupsReturn
}

func (v emptyNameValidator) ValidateModCliNames(namesToVerify []string, cliMods mc.ModMap) error {
	return v.ModsReturn
}

type vNameValidator struct {
	emptyNameValidator
	Groups       []string
	Mods         []string
	Map          mc.ModMap
	VisitedGroup *bool
	VisitedMod   *bool
}

func (v vNameValidator) ValidateServerGroups(groups []string) error {
	*(v.VisitedGroup) = true
	Expect(groups).To(ConsistOf(v.Groups))
	return v.GroupsReturn
}

func (v vNameValidator) ValidateModCliNames(namesToVerify []string, cliMods mc.ModMap) error {
	*(v.VisitedMod) = true
	Expect(namesToVerify).To(ConsistOf(v.Mods))
	Expect(cliMods).To(Equal(v.Map))
	return v.ModsReturn
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

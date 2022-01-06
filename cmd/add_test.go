package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"mcmods/cmd"
	"mcmods/mc"
	"strings"

	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Add Cmd", func() {
	var fs mc.FileSystem

	var mapValidator *nameMapperValidator
	var cfgIoSpy *clientConfigIoSpy
	var friendlyNoOp *noOpPrompt
	var cliNameNoOp *noOpPrompt
	var descNoOp *noOpPrompt
	var detailURLNoOp *noOpPrompt
	var latestURLNoOp *noOpPrompt
	var groupNoOp *noOpPrompt
	var serverAddSaveFake *serverAddSaveNoOp

	groupName := "optional"

	BeforeEach(func() {
		InitTestData()
		afs := afero.NewMemMapFs()
		fs = &mc.LocalFileSystem{Fs: afs}
		cmd.ViperInstance.SetFs(afs)

		mapb := false
		mapValidator = &nameMapperValidator{
			ClientMods: TestingClientMods,
			Visited:    &mapb,
			fakeNameMapper: fakeNameMapper{
				Map: TestingCliModMap,
			},
		}

		cfg := TestingConfig
		cfg.ClientMods = TestingClientMods
		b := false

		cfgIoSpy = &clientConfigIoSpy{
			Saved:      &b,
			LoadReturn: cfg,
		}

		serverAddSaveFake = &serverAddSaveNoOp{}

		cmd.NameMapper = mapValidator
		cmd.ServerCfgSaver = serverAddSaveFake
		cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
			return cfgIoSpy
		}
		cmd.CreateFsFunc = func(f *mc.FtpArgs) (mc.FileSystem, error) {
			return fs, nil
		}

		friendlyNoOp = &noOpPrompt{}
		cliNameNoOp = &noOpPrompt{}
		descNoOp = &noOpPrompt{}
		detailURLNoOp = &noOpPrompt{}
		latestURLNoOp = &noOpPrompt{}
		groupNoOp = &noOpPrompt{ReturnStr: groupName}

		cmd.FriendlyPrompt = friendlyNoOp
		cmd.CliNamePrompt = cliNameNoOp
		cmd.DescPrompt = descNoOp
		cmd.DetailsUrlPrompt = detailURLNoOp
		cmd.DownloadUrlPrompt = latestURLNoOp
		cmd.GroupPrompt = groupNoOp

		mc.ServerGroups = TestingServerGroups

		cmd.ResetAddVars()

		cmd.RootCmd.SetArgs([]string{"add"})
	})

	Context("errors", func() {
		expectedErr := errors.New("add error")

		It("returns error from friendly name prompt", func() {
			friendlyNoOp.ReturnErr = expectedErr

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from cli name prompt", func() {
			cliNameNoOp.ReturnErr = expectedErr

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from description prompt", func() {
			descNoOp.ReturnErr = expectedErr

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from detail URL prompt", func() {
			detailURLNoOp.ReturnErr = expectedErr

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from download url prompt", func() {
			latestURLNoOp.ReturnErr = expectedErr

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from group prompt", func() {
			groupNoOp.ReturnErr = expectedErr
			cmd.RootCmd.SetArgs([]string{"add", "--server"})

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})

		It("returns error from config saving", func() {
			serverAddSaveFake.Return = expectedErr
			cmd.RootCmd.SetArgs([]string{"add", "--server"})

			err := cmd.RootCmd.Execute()

			Expect(err).To(Equal(expectedErr))
		})
	})

	When("all valid inputs", func() {
		var expectedModValues *mc.Mod

		BeforeEach(func() {
			expectedModValues = &mc.Mod{
				FriendlyName: "Some fun mod",
				CliName:      "funmod",
				Description:  "A super duper fun mod",
				DetailsURL:   "<pretend this is a url>",
				LatestURL:    "<also pretend this is a url>",
			}

			friendlyNoOp.ReturnStr = expectedModValues.FriendlyName
			cliNameNoOp.ReturnStr = expectedModValues.CliName
			descNoOp.ReturnStr = expectedModValues.Description
			detailURLNoOp.ReturnStr = expectedModValues.DetailsURL
			latestURLNoOp.ReturnStr = expectedModValues.LatestURL
		})

		Context("client mod", func() {
			var clientAddIo *clientAddIoValidator

			BeforeEach(func() {
				b := false
				clientAddIo = &clientAddIoValidator{
					LoadReturn:      TestingConfig,
					ExpectedModCopy: expectedModValues,
					Saved:           &b,
				}

				cmd.ConfigIoFunc = func(f mc.FileSystem) mc.ModConfigIo {
					return clientAddIo
				}
			})

			It("adds a new mod to the client install config before saving", func() {
				cmd.RootCmd.Execute()

				// The clientAddIo validator ensures the item was added properly
				// This just makes sure that the validator was called
				Expect(*clientAddIo.Saved).To(BeTrue())
			})
		})

		Context("server mod", func() {
			var serverAddSave *serverAddSaveValidator

			BeforeEach(func() {
				b := false
				serverAddSave = &serverAddSaveValidator{
					ExpectedGroup:   groupName,
					ExpectedModCopy: expectedModValues,
					Saved:           &b,
				}

				cmd.ServerCfgSaver = serverAddSave

				cmd.RootCmd.SetArgs([]string{"add", "--server"})
			})

			It("adds a new mod to the server config before saving", func() {
				cmd.RootCmd.Execute()

				// The serverAddSave validator ensures the item was added properly
				// This just makes sure that the validator was called
				Expect(*serverAddSave.Saved).To(BeTrue())
			})
		})
	})

	Describe("prompt logic", func() {
		var inBuffer *bytes.Buffer
		var outBuffer *bytes.Buffer

		BeforeEach(func() {
			cmd.InitPrompts()
			outBuffer = bytes.NewBufferString("")
			inBuffer = bytes.NewBufferString("")
		})

		Describe("CLI name prompt", func() {
			It("allows valid names", func() {
				validNames := []string{
					"ab", "a-b", "a-b-c", "testname", "unreasonably-long-but-still-valid-name",
				}

				for _, name := range validNames {
					inBuffer.WriteString(name + "\n")
					str, err := cmd.CliNamePrompt.GetInput(outBuffer, inBuffer)

					Expect(err).To(BeNil())
					Expect(str).To(Equal(name))
					Expect(*mapValidator.Visited).To(BeTrue())
				}
			})

			It("rejects invalid names", func() {
				validName := "aaa" // last item must be valid to end the prompt loop
				invalidNames := strings.Join([]string{"am1", "a-b-", "-a-b-c", "TestName", "name2",
					"mod1.2.3", "mod_name", "mod+name", "mod name", "mod@name", "mod/name",
				}, "\n")
				inBuffer.WriteString(invalidNames + "\n" + validName + "\n")

				str, err := cmd.CliNamePrompt.GetInput(outBuffer, inBuffer)

				Expect(err).To(BeNil())
				Expect(str).To(Equal(validName))
			})
		})

		Describe("Server Group prompt", func() {
			It("allows server groups", func() {
				validNames := []string{
					"required", "optional", "performance", "server-only",
				}

				for _, name := range validNames {
					inBuffer.WriteString(name + "\n")
					str, err := cmd.GroupPrompt.GetInput(outBuffer, inBuffer)

					Expect(err).To(BeNil())
					Expect(str).To(Equal(name))
				}
			})

			It("rejects invalid groups", func() {
				inBuffer.WriteString("invalid\nrequired")

				str, err := cmd.GroupPrompt.GetInput(outBuffer, inBuffer)

				Expect(err).To(BeNil())
				Expect(str).To(Equal("required"))
			})
		})
	})
})

type noOpPrompt struct {
	ReturnStr string
	ReturnErr error
}

func (p noOpPrompt) GetInput(w io.Writer, r io.Reader) (string, error) {
	return p.ReturnStr, p.ReturnErr
}

type clientAddIoValidator struct {
	LoadReturn      *mc.UserModConfig
	ExpectedModCopy *mc.Mod
	Saved           *bool
}

func (i clientAddIoValidator) LoadOrNew() (*mc.UserModConfig, error) {
	return i.LoadReturn, nil
}

func (i clientAddIoValidator) Save(cfg *mc.UserModConfig) error {
	*(i.Saved) = true
	var mod, cMod *mc.Mod

	for _, cMod = range cfg.ClientMods {
		if cMod.CliName == i.ExpectedModCopy.CliName {
			mod = cMod
			break
		}
	}

	validateMod(mod, i.ExpectedModCopy)

	return nil
}

func validateMod(actual *mc.Mod, expected *mc.Mod) {
	Expect(actual).ToNot(BeNil())
	Expect(actual.FriendlyName).To(Equal(expected.FriendlyName))
	Expect(actual.Description).To(Equal(expected.Description))
	Expect(actual.DetailsURL).To(Equal(expected.DetailsURL))
	Expect(actual.LatestURL).To(Equal(expected.LatestURL))
}

type serverAddSaveNoOp struct {
	Return error
}

func (s serverAddSaveNoOp) Save() error {
	return s.Return
}

type serverAddSaveValidator struct {
	Return          error
	ExpectedModCopy *mc.Mod
	ExpectedGroup   string
	Saved           *bool
}

func (v serverAddSaveValidator) Save() error {
	*(v.Saved) = true
	var mod, cMod *mc.Mod
	var index int
	foundInGroups := []struct {
		name  string
		index int
	}{}

	for name, group := range mc.ServerGroups {
		for index, cMod = range group.Mods {
			if cMod.CliName == v.ExpectedModCopy.CliName {
				mod = cMod
				foundInGroups = append(foundInGroups, struct {
					name  string
					index int
				}{
					name:  name,
					index: index,
				})
				break
			}
		}
	}

	Expect(foundInGroups).To(HaveLen(1))
	Expect(foundInGroups[0].name).To(Equal(v.ExpectedGroup))

	validateMod(mod, v.ExpectedModCopy)

	return nil
}

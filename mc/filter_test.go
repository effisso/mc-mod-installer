package mc_test

import (
	"errors"
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var filter mc.ModFilter
	var cfg *mc.ClientModConfig
	var emptyMapper *emptyNameMapper
	var emptyValidator *emptyNameValidator

	xEmpty := []string{}

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		emptyValidator = &emptyNameValidator{}
		emptyMapper = &emptyNameMapper{}
		filter = mc.NewModFilter(emptyMapper, emptyValidator)
		cfg = &mc.ClientModConfig{
			ModInstallations: map[string]mc.ModInstallation{},
			ClientMods:       TestingClientMods,
		}
	})

	Context("exclusions", func() {
		When("no exclusions", func() {
			It("returns all mods", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, cfg, false)

				Expect(err).To(BeNil())
				Expect(mods).To(ConsistOf(TestingAllMods))
			})
		})

		When("server groups excluded", func() {
			It("excludes the specified mods", func() {
				mods, err := filter.FilterAllMods([]string{"performance", "optional"}, xEmpty, cfg, false)

				Expect(err).To(BeNil())
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})

		When("server mods excluded", func() {
			It("excludes the specified mods", func() {
				xServer := []string{TestingServerPerformance1.CliName, TestingServerOptional1.CliName}
				mods, err := filter.FilterAllMods(xEmpty, xServer, cfg, false)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})

		When("client mods excluded", func() {
			It("excludes the specified mods", func() {
				xClientMods := []string{TestingClientMod1.CliName, TestingClientMod2.CliName}
				mods, err := filter.FilterAllMods(xEmpty, xClientMods, cfg, false)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingClientMod1))
				Expect(mods).NotTo(ContainElement(TestingClientMod2))
			})
		})
	})

	Context("force", func() {
		BeforeEach(func() {
			install := mc.ModInstallation{DownloadUrl: TestingClientMod1.LatestUrl}
			cfg.ModInstallations[TestingClientMod1.CliName] = install
		})
		When("false", func() {
			It("doesn't return items with installations from the latest url", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, cfg, false)

				Expect(err).To(BeNil())
				Expect(mods).NotTo(ContainElement(TestingClientMod1))
			})
		})
		When("true", func() {
			It("returns items with installations from the latest url", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, cfg, true)

				Expect(err).To(BeNil())
				Expect(mods).To(ContainElement(TestingClientMod1))
			})

			It("does not return excluded items", func() {
				xServer := []string{TestingServerPerformance1.CliName, TestingServerOptional1.CliName}

				mods, err := filter.FilterAllMods(xEmpty, xServer, cfg, true)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})
	})

	Context("name verification", func() {
		Context("handles arguments correctly", func() {
			var xServer, xClient []string
			var nameMapper nameMapperValidator

			BeforeEach(func() {
				xServer = []string{TestingServerPerformance1.CliName, TestingServerOptional1.CliName}
				xClient = []string{TestingClientMod1.CliName}
				mapb := false
				nameMapper = nameMapperValidator{
					ClientMods:      TestingConfig.ClientMods,
					emptyNameMapper: emptyNameMapper{Map: TestingCliModMap},
					Visited:         &mapb,
				}
			})

			It("validates names", func() {
				groupb := false
				modb := false
				nmValidator := vNameValidator{
					Groups:             xServer,
					Mods:               xClient,
					Map:                TestingCliModMap,
					VisitedGroup:       &groupb,
					VisitedMod:         &modb,
					emptyNameValidator: emptyNameValidator{},
				}
				filter = mc.NewModFilter(nameMapper, nmValidator)

				mods, err := filter.FilterAllMods(xServer, xClient, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(Not(BeNil()))
				Expect(*nmValidator.VisitedGroup).To(BeTrue())
				Expect(*nameMapper.Visited).To(BeTrue())
				Expect(*nmValidator.VisitedMod).To(BeTrue())
			})
		})

		Context("on error", func() {
			It("returns server groups validation error", func() {
				(*emptyValidator).GroupsReturn = errors.New("group validation problem")

				mods, err := filter.FilterAllMods([]string{"invalid"}, xEmpty, cfg, true)

				Expect(err).To(Equal(emptyValidator.GroupsReturn))
				Expect(mods).To(BeNil())
			})

			It("returns mod name validation error", func() {
				emptyValidator.ModsReturn = errors.New("mod validation problem")

				mods, err := filter.FilterAllMods(xEmpty, []string{"invalid"}, cfg, true)

				Expect(err).To(Equal(emptyValidator.ModsReturn))
				Expect(mods).To(BeNil())
			})
		})
	})
})

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

type emptyNameMapper struct {
	Map mc.ModMap
}

func (m emptyNameMapper) MapAllMods(clientMods []*mc.Mod) mc.ModMap {
	return m.Map
}

type nameMapperValidator struct {
	emptyNameMapper
	ClientMods []*mc.Mod
	Visited    *bool
}

func (m nameMapperValidator) MapAllMods(clientMods []*mc.Mod) mc.ModMap {
	*(m.Visited) = true
	Expect(clientMods).To(ConsistOf(m.ClientMods))
	return m.Map
}

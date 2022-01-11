package mc_test

import (
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var filter mc.ModFilter
	var emptyMapper *emptyNameMapper

	xEmpty := []string{}

	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
		emptyMapper = &emptyNameMapper{Map: TestingCliModMap}
		filter = mc.NewModFilter(emptyMapper)
	})

	Context("exclusions", func() {
		BeforeEach(func() {
			// clear installed mods
			TestingConfig.ModInstallations = map[string]mc.ModInstallation{}
		})

		When("no exclusions", func() {
			It("returns all mods", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(ConsistOf(TestingAllMods))
			})
		})

		When("server groups excluded", func() {
			It("excludes the specified mods", func() {
				mods, err := filter.FilterAllMods([]string{"performance", "optional"}, xEmpty, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})

		When("server mods excluded", func() {
			It("excludes the specified mods", func() {
				xServer := []string{TestingServerPerformance1.CliName, TestingServerOptional1.CliName}

				mods, err := filter.FilterAllMods(xEmpty, xServer, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})

		When("client mods excluded", func() {
			It("excludes the specified mods", func() {
				xClientMods := []string{TestingClientMod1.CliName, TestingClientMod2.CliName}

				mods, err := filter.FilterAllMods(xEmpty, xClientMods, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingClientMod1))
				Expect(mods).NotTo(ContainElement(TestingClientMod2))
			})
		})

		When("validating server group names", func() {
			It("should return no error for existing server group names", func() {
				xServerGroups := []string{"performance"}

				mods, err := filter.FilterAllMods(xServerGroups, xEmpty, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(ConsistOf(
					// no performance
					TestingServerOnly1,
					TestingServerRequired1,
					TestingServerOptional1,
					TestingClientMod1,
					TestingClientMod2,
				))
			})

			It("should return an error for unknown server group names", func() {
				xServerGroups := []string{"invalid"}

				_, err := filter.FilterAllMods(xServerGroups, xEmpty, TestingConfig, false)

				Expect(err).ToNot(BeNil())
			})
		})

		When("validating mod CLI names", func() {
			It("should return no error for existing CLI names", func() {
				xClientMods := []string{TestingClientMod2.CliName}

				mods, err := filter.FilterAllMods(xEmpty, xClientMods, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).To(ConsistOf(
					// no client mod 2
					TestingServerOnly1,
					TestingServerRequired1,
					TestingServerOptional1,
					TestingServerPerformance1,
					TestingClientMod1,
				))
			})

			It("should return an error for unknown CLI names", func() {
				xClientMods := []string{"invalid"}

				_, err := filter.FilterAllMods(xEmpty, xClientMods, TestingConfig, false)

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("force", func() {
		BeforeEach(func() {
			install := mc.ModInstallation{DownloadURL: TestingClientMod1.LatestURL}
			TestingConfig.ModInstallations[TestingClientMod1.CliName] = install
		})
		When("false", func() {
			It("doesn't return items with installations from the latest url", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, TestingConfig, false)

				Expect(err).To(BeNil())
				Expect(mods).NotTo(ContainElement(TestingClientMod1))
			})
		})
		When("true", func() {
			It("returns items with installations from the latest url", func() {
				mods, err := filter.FilterAllMods(xEmpty, xEmpty, TestingConfig, true)

				Expect(err).To(BeNil())
				Expect(mods).To(ContainElement(TestingClientMod1))
			})

			It("does not return excluded items", func() {
				xServer := []string{"performance", "optional"}

				mods, err := filter.FilterAllMods(xServer, xEmpty, TestingConfig, true)

				Expect(err).To(BeNil())
				Expect(mods).To(HaveLen(len(TestingAllMods) - 2))
				Expect(mods).NotTo(ContainElement(TestingServerPerformance1))
				Expect(mods).NotTo(ContainElement(TestingServerOptional1))
			})
		})
	})
})

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

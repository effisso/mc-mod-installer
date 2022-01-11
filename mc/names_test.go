package mc_test

import (
	"mcmods/mc"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI", func() {
	BeforeEach(func() {
		InitTestData()
		mc.ServerGroups = TestingServerGroups
	})

	Describe("Mod Name Mapper", func() {
		var mapper mc.ModNameMapper

		BeforeEach(func() {
			mapper = mc.NewModNameMapper()
		})

		It("should return a map with all server and client mods, mapped by their CLI name", func() {
			modMap := mapper.MapAllMods(TestingClientMods)

			Expect(modMap).To(HaveLen(len(TestingAllMods)))
			for _, m := range TestingAllMods {
				Expect(modMap).Should(HaveKeyWithValue(m.CliName, m))
			}
		})

		It("should panic if mod names collide during construction", func() {
			defer func() {
				if r := recover(); r == nil {
					Fail("Did not panic with duplicate names")
				}
			}()

			badClientMods := TestingClientMods
			badClientMods = append(badClientMods, TestingClientMods[0]) // add a mod to the list again
			mapper.MapAllMods(badClientMods)

			Fail("Did not panic with duplicate names")
		})
	})
})

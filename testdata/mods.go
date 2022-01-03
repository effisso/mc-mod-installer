package testdata

import (
	"mcmods/mc"
)

var (
	TestingClientMod1 *mc.Mod
	TestingClientMod2 *mc.Mod
	TestingClientMods []*mc.Mod

	TestingServerRequired1    *mc.Mod
	TestingServerPerformance1 *mc.Mod
	TestingServerOptional1    *mc.Mod
	TestingServerOnly1        *mc.Mod

	TestingGroupRequired    *mc.ServerGroup
	TestingGroupOptional    *mc.ServerGroup
	TestingGroupPerformance *mc.ServerGroup
	TestingGroupServerOnly  *mc.ServerGroup

	TestingServerGroups map[string]*mc.ServerGroup

	TestingServerGroupNames []string

	TestingServerMods []*mc.Mod

	TestingAllMods []*mc.Mod

	TestingCliModMap mc.ModMap

	TestingClientCliNames []string
	TestingServerCliNames []string

	TestingConfig *mc.ClientModConfig
)

func InitTestData() {
	TestingClientMod1 = &mc.Mod{
		FriendlyName: "client mod 1",
		CliName:      "mod1",
		Description:  "mod 1 description",
		DetailsUrl:   "https://mod_1_dot_com",
		LatestUrl:    "https://mod_1_dot_com/latest",
	}
	TestingClientMod2 = &mc.Mod{
		FriendlyName: "client mod #2",
		CliName:      "modtwo",
		Description:  "description of the great mod 2",
		DetailsUrl:   "https://second_mod_dot_gov",
		LatestUrl:    "https://second_mod_dot_gov/latest",
	}
	TestingClientMods = []*mc.Mod{TestingClientMod1, TestingClientMod2}

	TestingServerRequired1 = &mc.Mod{
		FriendlyName: "REQUIRED mod",
		CliName:      "required1",
		Description:  "REQUIRED mod description",
		DetailsUrl:   "https://som_mod_site/mod-name",
		LatestUrl:    "https://some_mod_site/download/123",
	}
	TestingServerPerformance1 = &mc.Mod{
		FriendlyName: "performance mod",
		CliName:      "perf1",
		Description:  "performance mod description",
		DetailsUrl:   "https://some_mod_site/another-mod-name",
		LatestUrl:    "https://some_mod_site/download/547",
	}
	TestingServerOptional1 = &mc.Mod{
		FriendlyName: "optional mod",
		CliName:      "opt1",
		Description:  "optional mod description",
		DetailsUrl:   "https://mod_site/a-mod",
		LatestUrl:    "https://mod_site/a-mod/a-mod-9.8.7",
	}
	TestingServerOnly1 = &mc.Mod{
		FriendlyName: "server-only mod",
		CliName:      "svr1",
		Description:  "server-only mod description",
		DetailsUrl:   "https://modzone/server-mod",
		LatestUrl:    "https://modzone/server-mod/owiefaijdfaj/492834",
	}

	TestingGroupRequired = &mc.ServerGroup{
		Description: "mods required on the client",
		Mods:        []*mc.Mod{TestingServerRequired1},
	}
	TestingGroupOptional = &mc.ServerGroup{
		Description: "mods optional on the client",
		Mods:        []*mc.Mod{TestingServerOptional1},
	}
	TestingGroupPerformance = &mc.ServerGroup{
		Description: "performance mods optional on the client",
		Mods:        []*mc.Mod{TestingServerPerformance1},
	}
	TestingGroupServerOnly = &mc.ServerGroup{
		Description: "server mods that should not be on the client",
		Mods:        []*mc.Mod{TestingServerOnly1},
	}

	TestingServerGroups = map[string]*mc.ServerGroup{
		"required":    TestingGroupRequired,
		"optional":    TestingGroupOptional,
		"performance": TestingGroupPerformance,
		"server-only": TestingGroupServerOnly,
	}

	TestingServerGroupNames = []string{
		"required",
		"optional",
		"performance",
		"server-only",
	}

	TestingServerMods = []*mc.Mod{
		TestingServerRequired1,
		TestingServerOptional1,
		TestingServerPerformance1,
		TestingServerOnly1,
	}

	TestingAllMods = []*mc.Mod{
		TestingClientMod1,
		TestingClientMod2,
		TestingServerRequired1,
		TestingServerOptional1,
		TestingServerPerformance1,
		TestingServerOnly1,
	}

	TestingCliModMap = mc.ModMap{
		TestingClientMod1.CliName:         TestingClientMod1,
		TestingClientMod2.CliName:         TestingClientMod2,
		TestingServerOnly1.CliName:        TestingServerOnly1,
		TestingServerRequired1.CliName:    TestingServerRequired1,
		TestingServerOptional1.CliName:    TestingServerOptional1,
		TestingServerPerformance1.CliName: TestingServerPerformance1,
	}

	TestingClientCliNames = []string{
		TestingClientMod1.CliName,
		TestingClientMod2.CliName,
	}

	TestingServerCliNames = []string{
		TestingServerOnly1.CliName,
		TestingServerOptional1.CliName,
		TestingServerPerformance1.CliName,
		TestingServerRequired1.CliName,
	}

	TestingConfig = &mc.ClientModConfig{
		ModInstallations: map[string]mc.ModInstallation{
			TestingClientMod1.CliName: mc.ModInstallation{
				DownloadUrl: "dummy_url",
				Timestamp:   "123",
			},
			TestingServerRequired1.CliName: mc.ModInstallation{
				DownloadUrl: TestingServerRequired1.LatestUrl,
				Timestamp:   "789",
			},
		},
		ClientMods: []*mc.Mod{
			TestingClientMod1,
			TestingClientMod2,
		},
	}
}

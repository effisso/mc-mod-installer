package mc

// ModFilter is responsible for filtering the list of all mods down to those the user wishes to install
type ModFilter interface {
	// GetModsToInstall filters out mods as indicated by the user
	FilterAllMods(xGroups []string, xMods []string, cfg *UserModConfig, force bool) ([]*Mod, error)
}

type modFilter struct {
	NameMapper    ModNameMapper
	NameValidator NameValidator
}

// NewModFilter returns a new instance which implements ModFilter
func NewModFilter(mapper ModNameMapper, validator NameValidator) ModFilter {
	return modFilter{
		NameMapper:    mapper,
		NameValidator: validator,
	}
}

// GetModsToInstall filters out mods as indicated by the user
func (f modFilter) FilterAllMods(xGroups []string, xMods []string, cfg *UserModConfig, force bool) ([]*Mod, error) {
	mods := []*Mod{}
	xGroupSet := toSet(xGroups)
	xModSet := toSet(xMods)

	err := f.NameValidator.ValidateServerGroups(xGroups)
	if err != nil {
		return nil, err
	}

	modMap := f.NameMapper.MapAllMods(cfg.ClientMods)
	err = f.NameValidator.ValidateModCliNames(xMods, modMap)
	if err != nil {
		return nil, err
	}

	for groupName, group := range ServerGroups {
		if _, exclude := xGroupSet[groupName]; exclude {
			continue
		}
		for _, mod := range group.Mods {
			_, exclude := xModSet[mod.CliName]
			if !exclude && (!latestInstalled(mod, cfg) || force) {
				mods = append(mods, mod)
			}
		}
	}

	for _, mod := range cfg.ClientMods {
		_, exclude := xModSet[mod.CliName]
		if exclude || latestInstalled(mod, cfg) && !force {
			continue
		}
		mods = append(mods, mod)
	}

	return mods, nil
}

func latestInstalled(mod *Mod, cfg *UserModConfig) bool {
	latestInstalled := false
	installation, exists := cfg.ModInstallations[mod.CliName]

	if exists {
		latestInstalled = mod.LatestURL == installation.DownloadURL
	}

	return latestInstalled
}

func toSet(s []string) map[string]bool {
	set := map[string]bool{}

	for _, str := range s {
		set[str] = true
	}

	return set
}

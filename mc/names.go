package mc

import (
	"fmt"
)

type ModMap map[string]*Mod

// NameValidator helps validate user-provided names for things within the CLI
type NameValidator interface {
	// ValidateServerGroups returns an error if the server group is not valid
	ValidateServerGroups(groups []string) error

	// ValidateModCliNames returns an error if any of the mod names are not valid
	ValidateModCliNames(namesToVerify []string, cliMods ModMap) error
}

type nameValidator struct{}

// NewNameValidator returns an instance which implements NameValidator
func NewNameValidator() NameValidator {
	return nameValidator{}
}

// ValidateServerGroups returns an error if the server group is not valid
func (nameValidator) ValidateServerGroups(groups []string) error {
	for _, group := range groups {
		if _, exists := ServerGroups[group]; !exists {
			return fmt.Errorf("Unknown Server Group: %s", group)
		}
	}

	return nil
}

// ValidateModCliNames returns an error if any of the mod names are not valid
func (nameValidator) ValidateModCliNames(namesToVerify []string, cliMods ModMap) error {
	for _, name := range namesToVerify {
		if _, exists := cliMods[name]; !exists {
			return fmt.Errorf("Unknown Mod: %s", name)
		}
	}

	return nil
}

// ModNameMapper creates a map of mod CLI names to the mod definition
type ModNameMapper interface {
	// GetModMap returns a map of both client and server mod CLI names keyed to their mod definition
	MapAllMods(clientMods []*Mod) ModMap
}

type modNameMapper struct{}

// NewModNameMapper returns an instance which implements ModNameMapper
func NewModNameMapper() ModNameMapper {
	return modNameMapper{}
}

// GetModMap returns a map of cli names keyed to their Jar definition
func (_ modNameMapper) MapAllMods(clientMods []*Mod) ModMap {
	validNames := ModMap{}

	for _, group := range ServerGroups {
		addCliNamesToMap(group.Mods, validNames)
	}

	addCliNamesToMap(clientMods, validNames)

	return validNames
}

func addCliNamesToMap(modSlice []*Mod, m ModMap) {
	for _, mod := range modSlice {
		if _, exists := m[mod.CliName]; exists {
			panic(fmt.Errorf("Invalid config state: multiple mods share the same CLI name: %s", mod.CliName))
		}
		m[mod.CliName] = mod
	}
}

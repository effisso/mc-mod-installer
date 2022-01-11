package mc

import (
	"fmt"
)

// ModMap is shorthand for the common mapping of mod CLI name to mod definition
type ModMap map[string]*Mod

// ModNameMapper creates a map of mod CLI names to the mod definition
type ModNameMapper interface {
	// GetModMap returns a map of both client and server mod CLI names keyed to their
	// mod definition
	MapAllMods(clientMods []*Mod) ModMap
}

type modNameMapper struct{}

// NewModNameMapper returns an instance which implements ModNameMapper
func NewModNameMapper() ModNameMapper {
	return modNameMapper{}
}

// GetModMap returns a map of cli names keyed to their Jar definition
func (modNameMapper) MapAllMods(clientMods []*Mod) ModMap {
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

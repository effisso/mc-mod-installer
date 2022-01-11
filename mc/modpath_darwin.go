//go:build darwin
// +build darwin

package mc

import "path/filepath"

// DefaultOsMinecraftDir is where Minecraft is expected to be installed
var DefaultOsMinecraftDir = filepath.Join("~", "Library", "Application Support", "minecraft")

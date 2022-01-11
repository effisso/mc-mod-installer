//go:build linux
// +build linux

package mc

import (
	"path/filepath"
)

// DefaultOsMinecraftDir is where Minecraft is expected to be installed
var DefaultOsMinecraftDir = filepath.Join("$HOME", ".minecraft")

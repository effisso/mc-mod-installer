//go:build windows
// +build windows

package mc

import "path/filepath"

// DefaultOsMinecraftDir is where Minecraft is expected to be installed
var DefaultOsMinecraftDir = filepath.Join("%APPDATA%", ".minecraft")

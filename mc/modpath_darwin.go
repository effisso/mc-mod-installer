//go:build darwin
// +build darwin

package mc

import "path/filepath"

var DefaultOsMinecraftDir = filepath.Join("~", "Library", "Application Support", "minecraft")

//go:build windows
// +build windows

package mc

import "path/filepath"

var DefaultOsMinecraftDir = filepath.Join("%APPDATA%", ".minecraft")

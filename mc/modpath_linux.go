//go:build linux
// +build linux

package mc

import (
	"path/filepath"
)

var DefaultOsMinecraftDir = filepath.Join("$HOME", ".minecraft")

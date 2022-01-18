//go:build darwin
// +build darwin

package mc

import (
	"os"
	"path/filepath"
)

func init() {
	cfgDir, _ := os.UserHomeDir()
	DefaultOsMinecraftDir = filepath.Join(cfgDir, "Library", "Application Support", "minecraft")
}

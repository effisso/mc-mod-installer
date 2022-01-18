//go:build linux
// +build linux

package mc

import (
	"os"
	"path/filepath"
)

func init() {
	cfgDir, _ := os.UserHomeDir()
	DefaultOsMinecraftDir = filepath.Join(cfgDir, ".minecraft")
}

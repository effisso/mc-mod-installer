//go:build windows
// +build windows

package mc

import (
	"os"
	"path/filepath"
)

func init() {
	cfgDir, _ := os.UserConfigDir()
	DefaultOsMinecraftDir = filepath.Join(cfgDir, ".minecraft")
}

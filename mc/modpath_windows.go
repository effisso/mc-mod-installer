// +build windows

package mc

import "path/filepath"

var MinecraftDir = filepath.Join("%APPDATA%", ".minecraft")

// +build linux

package browser

import "os/exec"

func getBrowserCmd(url string) *exec.Cmd {
	return exec.Command("xdg-open", url)
}

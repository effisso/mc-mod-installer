// +build windows

package browser

import "os/exec"

func getBrowserCmd(url string) *exec.Cmd {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
}

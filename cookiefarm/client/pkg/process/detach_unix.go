//go:build !windows

package process

import (
	"os/exec"
	"syscall"
)

func setupDetach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
}

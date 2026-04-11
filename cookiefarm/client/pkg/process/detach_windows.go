//go:build windows

package process

import (
	"os/exec"
	"syscall"
)

func setupDetach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x08, // DETACHED_PROCESS
	}
}

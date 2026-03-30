//go:build !windows

package process

import (
	"os/exec"
	"syscall"
)

func setupCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func killProcess(cmd *exec.Cmd) error {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return err
	}
	return syscall.Kill(-pgid, syscall.SIGKILL)
}

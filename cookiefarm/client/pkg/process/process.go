package process

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

type Process struct {
	PID    int
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	cmd    *exec.Cmd
}

func Start(command string) (*Process, error) {
	return StartWithContext(context.Background(), command)
}

func StartDetached(command string) (*Process, error) {
	ctx := context.Background()
	cmd := buildCommand(ctx, command)
	setupCmd(cmd)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Process{
		PID:    cmd.Process.Pid,
		Stdout: nil,
		Stderr: nil,
		cmd:    cmd,
	}, nil
}

func StartWithContext(ctx context.Context, command string) (*Process, error) {
	cmd := buildCommand(ctx, command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe error: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe error: %w", err)
	}

	setupCmd(cmd)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start failed: %w", err)
	}

	return &Process{
		PID:    cmd.Process.Pid,
		Stdout: stdout,
		Stderr: stderr,
		cmd:    cmd,
	}, nil
}

func (p *Process) Wait() error {
	if p.cmd == nil {
		return fmt.Errorf("invalid process")
	}
	return p.cmd.Wait()
}

func (p *Process) Kill() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	return killProcess(p.cmd)
}

func buildCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd.exe", "/C", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

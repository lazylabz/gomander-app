//go:build !windows

package runner

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func newCommandProcess(command, workingDirectory string, environment []string) commandProcess {
	cmd := exec.Command(shellExecutable(), "-c", command)
	cmd.Dir = workingDirectory
	cmd.Env = environment
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return &execProcess{cmd: cmd}
}

func shellExecutable() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "/bin/sh"
	}
	return shell
}

func shouldSkipProcessOutputLine(string) bool {
	return false
}

func stopProcessGracefully(process commandProcess, done <-chan struct{}) error {
	select {
	case <-done:
		return nil
	default:
	}

	if err := syscall.Kill(-process.PID(), syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return ignoreNoSuchProcess(syscall.Kill(-process.PID(), syscall.SIGKILL))
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		return ignoreNoSuchProcess(syscall.Kill(-process.PID(), syscall.SIGKILL))
	case <-done:
		return nil
	}
}

func ignoreNoSuchProcess(err error) error {
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}

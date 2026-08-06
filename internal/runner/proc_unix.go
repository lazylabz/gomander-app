//go:build !windows

package runner

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

func newCommandProcess(command, workingDirectory string, environment []string) commandProcess {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", command)
	cmd.Dir = workingDirectory
	cmd.Env = environment
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return &execProcess{cmd: cmd}
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

	err := syscall.Kill(-process.PID(), syscall.SIGTERM)
	if err != nil {
		return syscall.Kill(-process.PID(), syscall.SIGKILL)
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		return syscall.Kill(-process.PID(), syscall.SIGKILL)
	case <-done:
		return nil
	}
}

//go:build !windows

package runner

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func newCommandProcess(command, workingDirectory string, environment []string, _ Config) commandProcess {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", command)
	cmd.Dir = workingDirectory
	cmd.Env = environment
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return &execProcess{cmd: cmd}
}

// stopProcessGracefully signals the process group and waits on exited, which
// the goroutine owning process.Wait closes.
func stopProcessGracefully(process commandProcess, exited <-chan struct{}) error {
	if alreadyExited(exited) {
		return nil
	}

	if err := syscall.Kill(-process.PID(), syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return forceKill(process)
	}

	select {
	case <-time.After(5 * time.Second):
		return forceKill(process)
	case <-exited:
		return nil
	}
}

func forceKill(process commandProcess) error {
	err := syscall.Kill(-process.PID(), syscall.SIGKILL)
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}

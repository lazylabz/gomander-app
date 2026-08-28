//go:build !windows

package runner

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func SetProcAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func SetProcEnv(cmd *exec.Cmd, environmentPaths []string) {
	if len(environmentPaths) == 0 {
		return
	}

	currentPath := os.Getenv("PATH")

	separator := ":"

	newPath := strings.Join(environmentPaths, separator) + separator + currentPath

	// Set the environment
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}

	// Update or add PATH
	for i, env := range cmd.Env {
		if strings.HasPrefix(strings.ToUpper(env), "PATH=") {
			cmd.Env[i] = "PATH=" + newPath
			return
		}
	}

	// If PATH wasn't found, add it
	cmd.Env = append(cmd.Env, "PATH="+newPath)
}

// StopProcessGracefully signals the process group and waits on exited, which the
// goroutine owning cmd.Wait closes. It must never call cmd.Wait itself.
func StopProcessGracefully(cmd *exec.Cmd, exited <-chan struct{}) error {
	if alreadyExited(exited) {
		return nil
	}

	// Try graceful termination first
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	if err != nil {
		if errors.Is(err, syscall.ESRCH) {
			// The Command ended on its own as it was being stopped
			return nil
		}
		// Fallback to force kill
		return forceKill(cmd)
	}

	select {
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown takes too long
		return forceKill(cmd)
	case <-exited:
		return nil
	}
}

func forceKill(cmd *exec.Cmd) error {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}

func GetCommand(cmdStr string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	return exec.Command(shell, "-c", cmdStr)
}

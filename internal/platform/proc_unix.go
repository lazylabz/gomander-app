//go:build !windows

package platform

import (
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

func SetProcEnv(cmd *exec.Cmd, extraPaths []string) {
	if len(extraPaths) == 0 {
		return
	}

	currentPath := os.Getenv("PATH")

	separator := ":"

	// Prepend extra paths to existing PATH
	newPath := strings.Join(extraPaths, separator) + separator + currentPath

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

func StopProcessGracefully(cmd *exec.Cmd) error {
	// Try graceful termination first
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	if err != nil {
		// Fallback to force kill
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	// Wait a bit for graceful shutdown
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown takes too long
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	case <-done:
		return nil
	}
}

func GetCommand(cmdStr string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	return exec.Command(shell, "-c", cmdStr)
}

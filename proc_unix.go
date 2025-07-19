//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

func setProcAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		Setpgid:    true,
	}
}

func stopProcessGracefully(cmd *exec.Cmd) error {
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

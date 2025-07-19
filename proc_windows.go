//go:build windows

package main

import (
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func setProcAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func stopProcessGracefully(cmd *exec.Cmd) error {
	pid := strconv.Itoa(cmd.Process.Pid)

	// Try graceful termination
	err := exec.Command("taskkill", "/PID", pid).Run()
	if err != nil {
		// Fallback to force kill
		return exec.Command("taskkill", "/F", "/T", "/PID", pid).Run()
	}

	// Wait for graceful shutdown
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// Force kill if needed
		return exec.Command("taskkill", "/F", "/T", "/PID", pid).Run()
	case <-done:
		return nil
	}
}

//go:build windows

package main

import (
	"os/exec"
	"strconv"
	"strings"
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
	killCmd := exec.Command("taskkill", "/PID", pid)
	killCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	err := killCmd.Run()
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

func getCommand(cmdStr string) *exec.Cmd {
	var cmd *exec.Cmd

	if strings.HasPrefix(cmdStr, "powershell ") {
		cmd = exec.Command("powershell", "-Command", strings.TrimPrefix(cmdStr, "powershell "))
	} else if strings.HasPrefix(cmdStr, "cmd ") {
		cmd = exec.Command("cmd", "/C", strings.TrimPrefix(cmdStr, "cmd "))
	} else {
		cmd = exec.Command("cmd", "/C", cmdStr)
	}

	return cmd
}

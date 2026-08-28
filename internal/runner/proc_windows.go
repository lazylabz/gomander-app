//go:build windows

package runner

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func SetProcAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func SetProcEnv(cmd *exec.Cmd, environmentPaths []string) {
	if len(environmentPaths) == 0 {
		return
	}

	currentPath := os.Getenv("PATH")

	separator := ";"

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

// StopProcessGracefully signals the process and waits on exited, which the
// goroutine owning cmd.Wait closes. It must never call cmd.Wait itself.
func StopProcessGracefully(cmd *exec.Cmd, exited <-chan struct{}) error {
	if alreadyExited(exited) {
		return nil
	}

	pid := strconv.Itoa(cmd.Process.Pid)

	// Try graceful termination
	killCmd := exec.Command("taskkill", "/PID", pid)
	killCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	err := killCmd.Run()
	if err != nil {
		// Fallback to force kill
		return forceKill(pid, exited)
	}

	select {
	case <-time.After(5 * time.Second):
		// Force kill if needed
		return forceKill(pid, exited)
	case <-exited:
		return nil
	}
}

// forceKill reports success when the Command is already gone: taskkill fails
// with "process not found" on a process that ended as it was being stopped.
func forceKill(pid string, exited <-chan struct{}) error {
	err := exec.Command("taskkill", "/F", "/T", "/PID", pid).Run()
	if err != nil && alreadyExited(exited) {
		return nil
	}
	return err
}

func GetCommand(cmdStr string) *exec.Cmd {
	cmd := exec.Command("cmd", "/C", cmdStr)

	return cmd
}

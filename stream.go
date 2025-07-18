package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	stdRuntime "runtime"
	"strings"
)

// ExecCommand executes a command by its ID and streams its output.
func (a *App) ExecCommand(id string) {
	command, exists := a.commands[id]
	if !exists {
		a.notifyError("Command not found: " + id)
		return
	}

	cmdStr := command.Command
	var cmd *exec.Cmd

	if stdRuntime.GOOS == "windows" {
		if strings.HasPrefix(cmdStr, "powershell ") {
			cmd = exec.Command("powershell", "-Command", strings.TrimPrefix(cmdStr, "powershell "))
		} else if strings.HasPrefix(cmdStr, "cmd ") {
			cmd = exec.Command("cmd", "/C", strings.TrimPrefix(cmdStr, "cmd "))
		} else {
			cmd = exec.Command("cmd", "/C", cmdStr)
		}
	} else {
		if strings.HasPrefix(cmdStr, "bash ") {
			cmd = exec.Command("bash", "-c", strings.TrimPrefix(cmdStr, "bash "))
		} else if strings.HasPrefix(cmdStr, "sh ") {
			cmd = exec.Command("sh", "-c", strings.TrimPrefix(cmdStr, "sh "))
		} else {
			cmd = exec.Command("sh", "-c", cmdStr)
		}
	}

	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.sendErrAsStreamLine(command, err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.sendErrAsStreamLine(command, err)
		return
	}

	if err := cmd.Start(); err != nil {
		a.sendErrAsStreamLine(command, err)
		return
	}
	a.commandsProcesses[command.Id] = cmd

	// Stream stdout
	go a.streamOutput(command.Id, stdout)
	// Stream stderr
	go a.streamOutput(command.Id, stderr)

	// Optional: Wait in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			a.logError(err.Error())
			return
		}
		a.emitEvent(ProcessFinished, command.Id)
	}()
}

func (a *App) streamOutput(commandId string, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)

	for scanner.Scan() {
		line := scanner.Text()
		a.logDebug(line)

		a.sendStreamLine(commandId, line)
	}
}

func (a *App) sendStreamLine(commandId string, line string) {
	a.emitEvent(NewLogEntry, map[string]string{
		"id":   commandId,
		"line": line,
	})
}

func (a *App) sendErrAsStreamLine(command Command, err error) {
	a.logDebug(err.Error())
	a.sendStreamLine(command.Id, err.Error())
}

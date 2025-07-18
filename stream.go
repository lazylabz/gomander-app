package main

import (
	"bufio"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os/exec"
	"strings"
)

// ExecCommand executes a command by its ID and streams its output.
func (a *App) ExecCommand(id string) {
	command, exists := a.commands[id]
	if !exists {
		a.notifyError("Command not found: " + id)
	}

	chunks := strings.Fields(command.Command)

	cmd := exec.Command(chunks[0], chunks[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.sendErrAsStreamLine(command, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.sendErrAsStreamLine(command, err)
	}

	if err := cmd.Start(); err != nil {
		a.sendErrAsStreamLine(command, err)
	}

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
		runtime.LogDebug(a.ctx, line)

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

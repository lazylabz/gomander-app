package main

import (
	"bufio"
	"io"
)

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

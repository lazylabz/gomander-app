package main

import (
	"bufio"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os/exec"
	"strings"
)

type Command struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
}

type LogServer struct {
	app      *App
	commands map[string]Command
}

func NewLogServer(app *App) *LogServer {
	return &LogServer{
		app:      app,
		commands: make(map[string]Command),
	}
}

func (l *LogServer) AddCommand(newCommand Command) {
	if _, exists := l.commands[newCommand.Id]; exists {
		l.app.logError("Command already exists: " + newCommand.Id)
		return
	}

	l.commands[newCommand.Id] = newCommand

	l.app.logInfo("Command added: " + newCommand.Id)

	l.app.emitEvent(GetCommands, nil)
}

func (l *LogServer) RemoveCommand(id string) {
	if _, exists := l.commands[id]; !exists {
		l.app.logError("Command not found: " + id)
		return
	}

	delete(l.commands, id)
	l.app.logInfo("Command removed: " + id)

	l.app.emitEvent(GetCommands, nil)
}

func (l *LogServer) EditCommand(newCommand Command) {
	if _, exists := l.commands[newCommand.Id]; !exists {
		l.app.logError("Command not found: " + newCommand.Id)
		return
	}

	l.commands[newCommand.Id] = newCommand
	l.app.logInfo("Command edited: " + newCommand.Id)

	l.app.emitEvent(GetCommands, nil)
}

func (l *LogServer) GetCommands() map[string]Command {
	return l.commands
}

func (l *LogServer) ExecCommand(id string) error {
	command, exists := l.commands[id]
	if !exists {
		l.app.logError("Command not found: " + id)
		return nil
	}

	chunks := strings.Fields(command.Command)

	cmd := exec.Command(chunks[0], chunks[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		l.app.logError(err.Error())
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		l.app.logError(err.Error())
		return err
	}

	if err := cmd.Start(); err != nil {
		l.app.logDebug(err.Error())
		return err
	}

	// Stream stdout
	go l.streamOutput(command.Id, stdout)
	// Stream stderr
	go l.streamOutput(command.Id, stderr)

	// Optional: Wait in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			l.app.logError(err.Error())
			return
		}
		l.app.emitEvent(ProcessFinished, command.Id)
	}()

	return nil
}

func (l *LogServer) streamOutput(commandId string, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)
	for scanner.Scan() {
		line := scanner.Text()
		runtime.LogDebug(l.app.ctx, line)

		l.app.emitEvent(NewLogEntry, map[string]string{
			"id":   commandId,
			"line": line,
		})
	}
}

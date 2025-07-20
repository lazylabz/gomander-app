package main

import (
	"os"
)

type Command struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
}

func (a *App) AddCommand(newCommand Command) {
	if _, exists := a.commands[newCommand.Id]; exists {
		a.logError("Command already exists: " + newCommand.Id)
		a.notifyError("Command already exists: " + newCommand.Id)
		return
	}

	a.commands[newCommand.Id] = newCommand

	a.logInfo("Command added: " + newCommand.Id)

	a.emitEvent(GetCommands, nil)

	err := saveConfig(&Config{
		Commands: a.commands,
	})

	if err != nil {
		a.notifyError(err.Error())
	}
}

func (a *App) RemoveCommand(id string) {
	if _, exists := a.commands[id]; !exists {
		a.logError("Command not found: " + id)
		a.notifyError("Command not found: " + id)
		return
	}

	delete(a.commands, id)
	a.logInfo("Command removed: " + id)

	a.emitEvent(GetCommands, nil)

	err := saveConfig(&Config{
		Commands: a.commands,
	})

	if err != nil {
		a.notifyError(err.Error())
	}
}

func (a *App) EditCommand(newCommand Command) {
	if _, exists := a.commands[newCommand.Id]; !exists {
		a.logError("Command not found: " + newCommand.Id)
		a.notifyError("Command not found: " + newCommand.Id)
		return
	}

	a.commands[newCommand.Id] = newCommand
	a.logInfo("Command edited: " + newCommand.Id)

	a.emitEvent(GetCommands, nil)

	err := saveConfig(&Config{
		Commands: a.commands,
	})

	if err != nil {
		a.notifyError(err.Error())
	}
}

func (a *App) GetCommands() map[string]Command {
	return a.commands
}

// ExecCommand executes a command by its ID and streams its output.
func (a *App) ExecCommand(id string) {
	command, exists := a.commands[id]
	if !exists {
		a.notifyError("Command not found: " + id)
		return
	}

	cmdStr := command.Command

	// Get the command object based on the command string and OS
	cmd := getCommand(cmdStr)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")

	// Set commabd attributes based on OS
	setProcAttributes(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.sendStreamError(command, err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.sendStreamError(command, err)
		return
	}

	if err := cmd.Start(); err != nil {
		a.sendStreamError(command, err)
		return
	}
	a.runningCommands[command.Id] = cmd

	// Stream stdout
	go a.streamOutput(command.Id, stdout)
	// Stream stderr
	go a.streamOutput(command.Id, stderr)

	// Optional: Wait in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			a.sendStreamError(command, err)
			a.logError("[ERROR - Waiting for command]: " + err.Error())
			return
		}
		a.emitEvent(ProcessFinished, command.Id)
	}()
}

func (a *App) StopRunningCommand(id string) error {
	cmd, exists := a.commands[id]

	if !exists {
		a.notifyError("Command not found: " + id)
		return nil
	}

	runningCommand, exists := a.runningCommands[cmd.Id]

	if !exists {
		a.notifyError("No running runningCommand for command: " + id)
		return nil
	}

	err := stopProcessGracefully(runningCommand)

	// If "graceful" stop fails, we try to kill the runningCommand
	if err != nil {
		a.logError(err.Error())
		a.notifyError("Failed to stop command gracefully: " + id + " - " + err.Error())
	}

	a.logInfo("Command stopped: " + id)

	a.emitEvent(ProcessFinished, cmd.Id)

	return nil
}

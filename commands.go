package main

import (
	"os"
	"os/exec"
	ntvRuntime "runtime"
	"strings"
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
	var cmd *exec.Cmd

	if ntvRuntime.GOOS == "windows" {
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
	a.commandsProcesses[command.Id] = cmd

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

	runningCommand, exists := a.commandsProcesses[cmd.Id]

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

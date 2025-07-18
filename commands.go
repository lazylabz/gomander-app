package main

import (
	ntvRuntime "runtime"
	"syscall"
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

func (a *App) StopRunningCommand(id string) error {
	cmd, exists := a.commands[id]

	if !exists {
		a.notifyError("Command not found: " + id)
		return nil
	}

	process, exists := a.commandsProcesses[cmd.Id]

	if !exists {
		a.notifyError("No running process for command: " + id)
		return nil
	}

	var err error

	if ntvRuntime.GOOS == "windows" {
		err = sendCtrlBreak(process.Process.Pid)
		if err != nil {
			a.logInfo(err.Error())
			// If sending CTRL_BREAK fails, try to kill the process directly
			err = process.Process.Kill()
			if err != nil {
				a.notifyError("Failed to stop command: " + id + " - " + err.Error())
				return err
			}
		}
	} else {
		err = process.Process.Signal(syscall.SIGTERM)
	}

	a.logInfo("Command stopped: " + id)

	a.emitEvent(ProcessFinished, cmd.Id)

	return nil
}

// TODO: Review why this only works sometimes
func sendCtrlBreak2(pid int) error {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return e
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return e
	}
	r, _, e := p.Call(uintptr(syscall.CTRL_BREAK_EVENT), uintptr(pid))
	if r == 0 {
		return e // syscall.GetLastError()
	}
	return nil
}

func sendCtrlBreak(pid int) error {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	generateConsoleCtrlEvent := kernel32.MustFindProc("GenerateConsoleCtrlEvent")

	// Enviar CTRL_BREAK al grupo de proceso (el PID debe ser el del proceso raíz creado con CREATE_NEW_PROCESS_GROUP)
	r, _, err := generateConsoleCtrlEvent.Call(uintptr(syscall.CTRL_BREAK_EVENT), uintptr(pid))
	if r == 0 {
		return err
	}
	return nil
}

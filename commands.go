package main

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

	err := process.Process.Kill()

	if err != nil {
		a.notifyError("Failed to stop command: " + id + " - " + err.Error())
		return err
	}

	a.logInfo("Command stopped: " + id)

	a.emitEvent(ProcessFinished, cmd.Id)
	
	return nil
}

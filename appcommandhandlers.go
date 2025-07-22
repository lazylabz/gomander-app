package main

func (a *App) GetCommands() map[string]Command {
	return a.commandsRepository.commands
}

func (a *App) AddCommand(newCommand Command) {
	err := a.commandsRepository.AddCommand(newCommand)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command added: " + newCommand.Id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command added")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.persistSavedConfig()
}

func (a *App) RemoveCommand(id string) {
	err := a.commandsRepository.RemoveCommand(id)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command removed: " + id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command removed")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.persistSavedConfig()
}

func (a *App) EditCommand(newCommand Command) {
	err := a.commandsRepository.EditCommand(newCommand)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command edited: " + newCommand.Id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command edited")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.persistSavedConfig()
}

func (a *App) RunCommand(id string) map[string]Command {
	command, err := a.commandsRepository.Get(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return nil
	}

	err = a.commandRunner.RunCommand(*command, a.savedConfig.ExtraPaths)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, "Failed to run command: "+id+" - "+err.Error())
		return nil
	}

	a.logger.info("Command executed: " + id)

	return nil
}

func (a *App) StopCommand(id string) {
	_, err := a.commandsRepository.Get(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	err = a.commandRunner.StopRunningCommand(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, "Failed to stop command gracefully: "+id+" - "+err.Error())
		return
	}

	a.logger.info("Command stopped: " + id)

	a.eventEmitter.emitEvent(ProcessFinished, id)
}

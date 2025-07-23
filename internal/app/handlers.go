package app

import (
	"gomander/internal/command"
	"gomander/internal/commandgroup"
	"gomander/internal/config"
	"gomander/internal/event"
)

// Command handlers

func (a *App) GetCommands() map[string]command.Command {
	return a.commandsRepository.GetCommands()
}

func (a *App) AddCommand(newCommand command.Command) {
	err := a.commandsRepository.AddCommand(newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return
	}

	a.logger.Info("Command added: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command added")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	a.persistRepositoryInformation()
}

func (a *App) RemoveCommand(id string) {
	err := a.commandsRepository.RemoveCommand(id)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return
	}

	a.commandGroupsRepository.RemoveCommandFromCommandGroups(id)

	a.logger.Info("Command removed: " + id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command removed")

	// Update the commands and command groups in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	a.persistRepositoryInformation()
}

func (a *App) EditCommand(newCommand command.Command) {
	err := a.commandsRepository.EditCommand(newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return
	}

	a.logger.Info("Command edited: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command edited")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	a.persistRepositoryInformation()
}

func (a *App) RunCommand(id string) map[string]command.Command {
	cmd, err := a.commandsRepository.GetCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return nil
	}

	extraPaths := a.extraPathRepository.GetExtraPaths()
	extraPathsStr := make([]string, len(extraPaths))
	for i, path := range extraPaths {
		extraPathsStr[i] = string(path)
	}

	err = a.commandRunner.RunCommand(*cmd, extraPathsStr)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to run command: "+id+" - "+err.Error())
		return nil
	}

	a.logger.Info("Command executed: " + id)

	return nil
}

func (a *App) StopCommand(id string) {
	_, err := a.commandsRepository.GetCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return
	}

	err = a.commandRunner.StopRunningCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to stop command gracefully: "+id+" - "+err.Error())
		return
	}

	a.logger.Info("Command stopped: " + id)

	a.eventEmitter.EmitEvent(event.ProcessFinished, id)
}

// Command group handlers

func (a *App) GetCommandGroups() []commandgroup.CommandGroup {
	return a.commandGroupsRepository.GetCommandGroups()
}

func (a *App) SaveCommandGroups(commandGroups []commandgroup.CommandGroup) {
	a.commandGroupsRepository.SetCommandGroups(commandGroups)
	a.persistRepositoryInformation()
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command groups saved successfully")
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)
}

// User config handlers

func (a *App) SaveUserConfig(userConfig config.UserConfig) {
	a.extraPathRepository.SetExtraPaths(userConfig.ExtraPaths)

	a.persistRepositoryInformation()

	a.logger.Info("Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.GetUserConfig, nil)
}

func (a *App) GetUserConfig() config.UserConfig {
	return config.UserConfig{
		ExtraPaths: a.extraPathRepository.GetExtraPaths(),
	}
}

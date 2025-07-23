package app

import (
	"gomander/internal/command"
	"gomander/internal/config"
	"gomander/internal/event"
	"slices"
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

	a.persistConfig()
}

func (a *App) RemoveCommand(id string) {
	err := a.commandsRepository.RemoveCommand(id)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return
	}

	a.savedConfig.CommandGroups = a.removeCommandFromGroups(id)

	a.logger.Info("Command removed: " + id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command removed")

	// Update the commands and command groups in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	a.persistConfig()
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

	a.persistConfig()
}

func (a *App) RunCommand(id string) map[string]command.Command {
	cmd, err := a.commandsRepository.GetCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return nil
	}

	err = a.commandRunner.RunCommand(*cmd, a.savedConfig.ExtraPaths)

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

func (a *App) GetCommandGroups() []config.CommandGroup {
	return a.savedConfig.CommandGroups
}

// Command group handlers

func (a *App) SaveCommandGroups(commandGroups []config.CommandGroup) {
	a.savedConfig.CommandGroups = commandGroups
	a.persistConfig()
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command groups saved successfully")
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)
}

func (a *App) removeCommandFromGroups(commandId string) []config.CommandGroup {
	newCommandGroups := make([]config.CommandGroup, 0)

	for _, group := range a.savedConfig.CommandGroups {
		if slices.Contains(group.CommandIds, commandId) {
			newCommandIds := make([]string, 0)

			for _, cmdId := range group.CommandIds {
				if cmdId != commandId {
					newCommandIds = append(newCommandIds, cmdId)
				}
			}

			group.CommandIds = newCommandIds
		}
		newCommandGroups = append(newCommandGroups, group)

	}

	return newCommandGroups
}

// User config handlers

func (a *App) SaveUserConfig(userConfig config.UserConfig) {
	a.savedConfig.ExtraPaths = userConfig.ExtraPaths

	a.persistConfig()

	a.logger.Info("Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.GetUserConfig, nil)
}

func (a *App) GetUserConfig() config.UserConfig {
	return config.UserConfig{
		ExtraPaths: a.savedConfig.ExtraPaths,
	}
}

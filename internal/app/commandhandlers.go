package app

import (
	"sort"

	"gomander/internal/command/domain"
	domain2 "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
)

func (a *App) GetCommands() ([]domain.Command, error) {
	return a.commandRepository.GetCommands(a.openedProjectId)
}

func (a *App) AddCommand(newCommand domain.Command) error {
	allCommands, err := a.commandRepository.GetCommands(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	newPosition := len(allCommands)

	newCommand.Position = newPosition

	err = a.commandRepository.SaveCommand(&newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	a.logger.Info("Command added: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command added")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	return nil
}

func (a *App) RemoveCommand(id string) error {
	err := a.commandRepository.DeleteCommand(id)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	a.RemoveCommandFromCommandGroups(id)

	a.logger.Info("Command removed: " + id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command removed")

	// Update the commands and command groups in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) EditCommand(newCommand domain.Command) error {
	err := a.commandRepository.EditCommand(&newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	a.logger.Info("Command edited: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command edited")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	return nil
}

func (a *App) ReorderCommands(orderedIds []string) error {
	existingCommands, err := a.commandRepository.GetCommands(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
	}
	// Sort the existing commands based on the new order
	sort.Slice(existingCommands, func(i, j int) bool {
		return array.IndexOf(orderedIds, existingCommands[i].Id) < array.IndexOf(orderedIds, existingCommands[j].Id)
	})

	// Update the position of each command based on the new order
	for i := range existingCommands {
		existingCommands[i].Position = i
		err := a.commandRepository.EditCommand(&existingCommands[i])
		if err != nil {
			a.logger.Error(err.Error())
			a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
			return err
		}
	}

	a.logger.Info("Commands reordered")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	return nil
}

func (a *App) RunCommand(id string) error {
	cmd, err := a.commandRepository.GetCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return nil
	}

	userConfig, err := a.userConfigRepository.GetOrCreateConfig()
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to get user config: "+err.Error())
		return nil
	}

	currentProject, err := a.projectRepository.GetProjectById(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to get current project: "+err.Error())
		return nil
	}

	environmentPathsStrings := array.Map(userConfig.EnvironmentPaths, func(ep domain2.EnvironmentPath) string { return ep.Path })
	err = a.commandRunner.RunCommand(*cmd, environmentPathsStrings, currentProject.WorkingDirectory)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to run command: "+id+" - "+err.Error())
		return nil
	}

	a.logger.Info("Command executed: " + id)

	return nil
}

func (a *App) StopCommand(id string) {
	// Check if the command exists before trying to stop it
	_, err := a.commandRepository.GetCommand(id)

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

package app

import (
	"gomander/internal/event"
	"gomander/internal/project"
)

func (a *App) GetCommands() map[string]project.Command {
	return a.selectedProject.GetCommands()
}

func (a *App) AddCommand(newCommand project.Command) error {
	err := a.selectedProject.AddCommand(newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	err = a.persistSelectedProjectConfig()
	if err != nil {
		return err
	}

	a.logger.Info("Command added: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command added")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	return nil
}

func (a *App) RemoveCommand(id string) error {
	err := a.selectedProject.RemoveCommand(id)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	a.selectedProject.RemoveCommandFromCommandGroups(id)

	err = a.persistSelectedProjectConfig()
	if err != nil {
		return err
	}

	a.logger.Info("Command removed: " + id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command removed")

	// Update the commands and command groups in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) EditCommand(newCommand project.Command) error {
	err := a.selectedProject.EditCommand(newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}

	err = a.persistSelectedProjectConfig()
	if err != nil {
		return err
	}

	a.logger.Info("Command edited: " + newCommand.Id)
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command edited")

	// Update the commands map in the frontend
	a.eventEmitter.EmitEvent(event.GetCommands, nil)

	return nil
}

func (a *App) RunCommand(id string) map[string]project.Command {
	cmd, err := a.selectedProject.GetCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return nil
	}

	extraPaths := a.userConfig.ExtraPaths
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
	_, err := a.selectedProject.GetCommand(id)

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

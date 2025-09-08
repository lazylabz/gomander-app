package app

import (
	domain2 "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
)

func (a *App) RunCommand(id string) error {
	cmd, err := a.commandRepository.Get(id)

	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	userConfig, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	currentProject, err := a.projectRepository.Get(userConfig.LastOpenedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	environmentPathsStrings := array.Map(userConfig.EnvironmentPaths, func(ep domain2.EnvironmentPath) string { return ep.Path })
	err = a.commandRunner.RunCommand(cmd, environmentPathsStrings, currentProject.WorkingDirectory)

	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	a.logger.Info("Command executed: " + id)

	return nil
}

func (a *App) StopCommand(id string) {
	// Check if the command exists before trying to stop it
	_, err := a.commandRepository.Get(id)

	if err != nil {
		a.logger.Error(err.Error())
		return
	}

	err = a.commandRunner.StopRunningCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		return
	}

	a.logger.Info("Command stopped: " + id)

	a.eventEmitter.EmitEvent(event.ProcessFinished, id)
}

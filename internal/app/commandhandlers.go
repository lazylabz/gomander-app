package app

import (
	"errors"
	"sort"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	domain2 "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
)

func (a *App) GetCommands() ([]domain.Command, error) {
	return a.commandRepository.GetAll(a.openedProjectId)
}

func (a *App) AddCommand(newCommand domain.Command) error {
	allCommands, err := a.commandRepository.GetAll(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	newPosition := len(allCommands)

	newCommand.Position = newPosition

	err = a.commandRepository.Create(&newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	a.logger.Info("Command added: " + newCommand.Id)

	return nil
}

func (a *App) RemoveCommand(id string) error {
	err := a.commandRepository.Delete(id)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	domainEvent := domainevent.NewCommandDeletedEvent(id)

	errs := a.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while removing command:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
			a.logger.Error(pubErr.Error())
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	a.logger.Info("Command removed: " + id)

	return nil
}

func (a *App) EditCommand(newCommand domain.Command) error {
	err := a.commandRepository.Update(&newCommand)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	a.logger.Info("Command edited: " + newCommand.Id)

	return nil
}

func (a *App) ReorderCommands(orderedIds []string) error {
	existingCommands, err := a.commandRepository.GetAll(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	// Sort the existing commands based on the new order
	sort.Slice(existingCommands, func(i, j int) bool {
		return array.IndexOf(orderedIds, existingCommands[i].Id) < array.IndexOf(orderedIds, existingCommands[j].Id)
	})

	// Update the position of each command based on the new order
	for i := range existingCommands {
		existingCommands[i].Position = i
		err := a.commandRepository.Update(&existingCommands[i])
		if err != nil {
			a.logger.Error(err.Error())
			return err
		}
	}

	a.logger.Info("Commands reordered")

	return nil
}

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

	currentProject, err := a.projectRepository.Get(a.openedProjectId)
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

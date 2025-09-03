package app

import (
	"errors"

	"gomander/internal/project/domain"
	"gomander/internal/project/domain/event"
)

func (a *App) OpenProject(projectConfigId string) error {
	_, err := a.projectRepository.Get(projectConfigId)
	if err != nil {
		return err
	}

	config, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectConfigId

	err = a.userConfigRepository.Update(config)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) CreateProject(project domain.Project) error {
	err := a.projectRepository.Create(project)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) EditProject(project domain.Project) error {
	err := a.projectRepository.Update(project)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) CloseProject() error {
	config, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = ""

	err = a.userConfigRepository.Update(config)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteProject(projectId string) error {
	err := a.projectRepository.Delete(projectId)
	if err != nil {
		return err
	}

	domainEvent := event.NewProjectDeletedEvent(projectId)

	errs := a.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while removing project:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
			a.logger.Error(pubErr.Error())
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	return nil
}

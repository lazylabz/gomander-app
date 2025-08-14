package app

import (
	"gomander/internal/event"
	"gomander/internal/project/domain"
)

func (a *App) GetCurrentProject() (*domain.Project, error) {
	return a.projectRepository.GetProjectById(a.openedProjectId)
}

func (a *App) GetAvailableProjects() ([]domain.Project, error) {
	return a.projectRepository.GetAllProjects()
}

func (a *App) OpenProject(projectConfigId string) error {
	config, err := a.userConfigRepository.GetOrCreateConfig()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectConfigId

	err = a.userConfigRepository.SaveConfig(config)
	if err != nil {
		return err
	}

	a.openedProjectId = projectConfigId

	return nil
}

func (a *App) CreateProject(project domain.Project) error {
	err := a.projectRepository.CreateProject(project)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) EditProject(project domain.Project) error {
	err := a.projectRepository.UpdateProject(project)
	if err != nil {
		return err
	}

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project edited successfully")

	return nil
}

func (a *App) CloseProject() error {
	config, err := a.userConfigRepository.GetOrCreateConfig()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = ""

	err = a.userConfigRepository.SaveConfig(config)
	if err != nil {
		return err
	}

	a.openedProjectId = ""

	return nil
}

func (a *App) DeleteProject(projectConfigId string) error {
	err := a.projectRepository.DeleteProject(projectConfigId)

	if err != nil {
		return err
	}

	return nil
}

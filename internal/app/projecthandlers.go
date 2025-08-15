package app

import (
	"gomander/internal/event"
	"gomander/internal/project/domain"
)

func (a *App) GetCurrentProject() (*domain.Project, error) {
	return a.projectRepository.Get(a.openedProjectId)
}

func (a *App) GetAvailableProjects() ([]domain.Project, error) {
	return a.projectRepository.GetAll()
}

func (a *App) OpenProject(projectConfigId string) error {
	config, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectConfigId

	err = a.userConfigRepository.Update(config)
	if err != nil {
		return err
	}

	a.openedProjectId = projectConfigId

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

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project edited successfully")

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

	a.openedProjectId = ""

	return nil
}

func (a *App) DeleteProject(projectConfigId string) error {
	commands, err := a.commandRepository.GetAll(projectConfigId)
	if err != nil {
		return err
	}

	for _, command := range commands {
		err := a.RemoveCommand(command.Id)
		if err != nil {
			return err
		}
	}

	err = a.projectRepository.Delete(projectConfigId)
	if err != nil {
		return err
	}

	return nil
}

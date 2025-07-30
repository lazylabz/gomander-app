package app

import (
	"gomander/internal/project"
)

func (a *App) GetCurrentProject() *project.Project {
	return a.selectedProject
}

func (a *App) GetAvailableProjects() ([]*project.Project, error) {
	return project.GetAllProjectsAvailableInProjectsFolder()
}

func (a *App) OpenProject(projectConfigId string) (*project.Project, error) {
	p, err := project.LoadProject(projectConfigId)
	if err != nil {
		return nil, err
	}

	a.userConfig.LastOpenedProjectId = p.Id

	err = a.persistUserConfig()
	if err != nil {
		return nil, err
	}

	a.selectedProject = p

	return p, nil
}

func (a *App) CreateProject(id, name string) error {
	err := project.SaveProject(&project.Project{
		Id:            id,
		Name:          name,
		Commands:      make(map[string]project.Command),
		CommandGroups: make([]project.CommandGroup, 0),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *App) CloseProject() error {
	if a.selectedProject == nil {
		return nil
	}

	a.selectedProject = nil

	a.userConfig.LastOpenedProjectId = ""

	err := a.persistUserConfig()
	if err != nil {
		return err
	}

	return nil
}

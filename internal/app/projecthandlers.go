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

	a.selectedProject = p

	return p, nil
}

package usecases

import (
	"gomander/internal/openedproject"
	"gomander/internal/project/domain"
)

type GetCurrentProject struct {
	openedProject openedproject.OpenedProject
}

func NewGetCurrentProject(openedProject openedproject.OpenedProject) *GetCurrentProject {
	return &GetCurrentProject{
		openedProject: openedProject,
	}
}

func (uc *GetCurrentProject) Execute() (*domain.Project, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return nil, err
	}
	if !open {
		return nil, nil
	}

	return &project, nil
}

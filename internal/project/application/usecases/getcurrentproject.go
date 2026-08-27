package usecases

import (
	"gomander/internal/openedproject"
	"gomander/internal/project/domain"
)

type GetCurrentProject interface {
	Execute() (*domain.Project, error)
}

type DefaultGetCurrentProject struct {
	openedProject openedproject.OpenedProject
}

func NewGetCurrentProject(openedProject openedproject.OpenedProject) *DefaultGetCurrentProject {
	return &DefaultGetCurrentProject{
		openedProject: openedProject,
	}
}

func (uc *DefaultGetCurrentProject) Execute() (*domain.Project, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return nil, err
	}
	if !open {
		return nil, nil
	}

	return &project, nil
}

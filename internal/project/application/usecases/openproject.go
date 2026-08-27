package usecases

import (
	"gomander/internal/openedproject"
)

type OpenProject interface {
	Execute(projectId string) error
}

type DefaultOpenProject struct {
	openedProject openedproject.OpenedProject
}

func NewOpenProject(openedProject openedproject.OpenedProject) *DefaultOpenProject {
	return &DefaultOpenProject{
		openedProject: openedProject,
	}
}

func (uc *DefaultOpenProject) Execute(projectId string) error {
	return uc.openedProject.Open(projectId)
}

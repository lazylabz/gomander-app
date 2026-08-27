package usecases

import (
	"gomander/internal/openedproject"
)

type CloseProject interface {
	Execute() error
}

type DefaultCloseProject struct {
	openedProject openedproject.OpenedProject
}

func NewCloseProject(openedProject openedproject.OpenedProject) *DefaultCloseProject {
	return &DefaultCloseProject{
		openedProject: openedProject,
	}
}

func (uc *DefaultCloseProject) Execute() error {
	return uc.openedProject.Close()
}

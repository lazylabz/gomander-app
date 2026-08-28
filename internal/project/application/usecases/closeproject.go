package usecases

import (
	"gomander/internal/openedproject"
)

type CloseProject struct {
	openedProject openedproject.OpenedProject
}

func NewCloseProject(openedProject openedproject.OpenedProject) *CloseProject {
	return &CloseProject{
		openedProject: openedProject,
	}
}

func (uc *CloseProject) Execute() error {
	return uc.openedProject.Close()
}

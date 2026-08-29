package usecases

import (
	"gomander/internal/openedproject"
)

type OpenProject struct {
	openedProject openedproject.OpenedProject
}

func NewOpenProject(openedProject openedproject.OpenedProject) *OpenProject {
	return &OpenProject{
		openedProject: openedProject,
	}
}

func (uc *OpenProject) Execute(projectId string) error {
	return uc.openedProject.Open(projectId)
}

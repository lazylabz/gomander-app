package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
)

type GetCommands struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewGetCommands(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *GetCommands {
	return &GetCommands{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *GetCommands) Execute() ([]domain.Command, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return make([]domain.Command, 0), err
	}
	if !open {
		return make([]domain.Command, 0), nil
	}

	return uc.commandRepository.GetAll(project.Id)
}

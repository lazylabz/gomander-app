package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
)

type GetCommands interface {
	Execute() ([]domain.Command, error)
}

type DefaultGetCommands struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewGetCommands(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *DefaultGetCommands {
	return &DefaultGetCommands{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultGetCommands) Execute() ([]domain.Command, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return make([]domain.Command, 0), err
	}
	if !open {
		return make([]domain.Command, 0), nil
	}

	return uc.commandRepository.GetAll(project.Id)
}

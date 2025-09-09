package usecases

import (
	"gomander/internal/commandgroup/domain"
)

type UpdateCommandGroup interface {
	Execute(commandGroup *domain.CommandGroup) error
}

type DefaultUpdateCommandGroup struct {
	commandGroupRepository domain.Repository
}

func NewUpdateCommandGroup(commandGroupRepo domain.Repository) *DefaultUpdateCommandGroup {
	return &DefaultUpdateCommandGroup{
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultUpdateCommandGroup) Execute(commandGroup *domain.CommandGroup) error {
	return uc.commandGroupRepository.Update(commandGroup)
}

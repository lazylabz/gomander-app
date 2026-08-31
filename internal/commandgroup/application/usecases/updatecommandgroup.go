package usecases

import (
	"gomander/internal/commandgroup/domain"
)

type UpdateCommandGroup struct {
	commandGroupRepository domain.Repository
}

func NewUpdateCommandGroup(commandGroupRepo domain.Repository) *UpdateCommandGroup {
	return &UpdateCommandGroup{
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *UpdateCommandGroup) Execute(commandGroup *domain.CommandGroup) error {
	return uc.commandGroupRepository.Update(commandGroup)
}

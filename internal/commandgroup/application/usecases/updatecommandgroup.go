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

func (uc *UpdateCommandGroup) Execute(commandGroup *domain.CommandGroupWithCommandIds) error {
	return uc.commandGroupRepository.UpdateWithCommandIds(commandGroup)
}

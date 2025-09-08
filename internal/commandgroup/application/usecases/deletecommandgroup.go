package usecases

import (
	"gomander/internal/commandgroup/domain"
)

type DeleteCommandGroup interface {
	Execute(commandGroupId string) error
}

type DefaultDeleteCommandGroup struct {
	commandGroupRepository domain.Repository
}

func NewDeleteCommandGroup(commandGroupRepo domain.Repository) *DefaultDeleteCommandGroup {
	return &DefaultDeleteCommandGroup{
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultDeleteCommandGroup) Execute(commandGroupId string) error {
	return uc.commandGroupRepository.Delete(commandGroupId)
}

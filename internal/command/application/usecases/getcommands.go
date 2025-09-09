package usecases

import (
	"gomander/internal/command/domain"
	configdomain "gomander/internal/config/domain"
)

type GetCommands interface {
	Execute() ([]domain.Command, error)
}

type DefaultGetCommands struct {
	configRepository  configdomain.Repository
	commandRepository domain.Repository
}

func NewGetCommands(configRepo configdomain.Repository, commandRepo domain.Repository) *DefaultGetCommands {
	return &DefaultGetCommands{
		configRepository:  configRepo,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultGetCommands) Execute() ([]domain.Command, error) {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return make([]domain.Command, 0), err
	}
	return uc.commandRepository.GetAll(userConfig.LastOpenedProjectId)
}

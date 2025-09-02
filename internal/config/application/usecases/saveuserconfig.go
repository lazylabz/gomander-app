package usecases

import (
	"gomander/internal/config/domain"
	"gomander/internal/logger"
)

type SaveUserConfig interface {
	Execute(config domain.Config) error
}

type DefaultSaveUserConfig struct {
	repository domain.Repository
	logger     logger.Logger
}

func NewSaveUserConfig(repository domain.Repository, logger logger.Logger) *DefaultSaveUserConfig {
	return &DefaultSaveUserConfig{repository: repository, logger: logger}
}

func (uc *DefaultSaveUserConfig) Execute(newUserConfig domain.Config) error {
	err := uc.repository.Update(&newUserConfig)
	if err != nil {
		uc.logger.Error("Failed to save user configuration: " + err.Error())
		return err
	}

	uc.logger.Info("User configuration saved successfully")
	return nil
}

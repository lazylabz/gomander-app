package usecases

import (
	"gomander/internal/config/domain"
)

type SaveUserConfig interface {
	Execute(config domain.Config) error
}

type DefaultSaveUserConfig struct {
	repository domain.Repository
}

func NewSaveUserConfig(repository domain.Repository) *DefaultSaveUserConfig {
	return &DefaultSaveUserConfig{repository: repository}
}

func (uc *DefaultSaveUserConfig) Execute(newUserConfig domain.Config) error {
	return uc.repository.Update(&newUserConfig)
}

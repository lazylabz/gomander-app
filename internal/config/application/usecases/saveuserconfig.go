package usecases

import (
	"gomander/internal/config/domain"
)

type SaveUserConfig struct {
	repository domain.Repository
}

func NewSaveUserConfig(repository domain.Repository) *SaveUserConfig {
	return &SaveUserConfig{repository: repository}
}

func (uc *SaveUserConfig) Execute(newUserConfig domain.Config) error {
	if err := newUserConfig.Validate(); err != nil {
		return err
	}

	return uc.repository.Update(&newUserConfig)
}

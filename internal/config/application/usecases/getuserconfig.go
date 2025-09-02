package usecases

import "gomander/internal/config/domain"

type GetUserConfig interface {
	Execute() (*domain.Config, error)
}

type DefaultGetUserConfig struct {
	repository domain.Repository
}

func NewDefaultGetUserConfig(repository domain.Repository) *DefaultGetUserConfig {
	return &DefaultGetUserConfig{repository: repository}
}

func (uc *DefaultGetUserConfig) Execute() (*domain.Config, error) {
	return uc.repository.GetOrCreate()
}

package usecases

import "gomander/internal/config/domain"

type GetUserConfig struct {
	repository domain.Repository
}

func NewGetUserConfig(repository domain.Repository) *GetUserConfig {
	return &GetUserConfig{repository: repository}
}

func (uc *GetUserConfig) Execute() (*domain.Config, error) {
	return uc.repository.GetOrCreate()
}

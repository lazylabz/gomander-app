package test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
)

type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockConfigRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

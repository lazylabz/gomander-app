package usecases_test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
)

type MockUserConfigRepository struct {
	mock.Mock
}

func (m *MockUserConfigRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockUserConfigRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
}

type MockEventEmitter struct {
	mock.Mock
}

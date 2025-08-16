package app_test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	"gomander/internal/event"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
	return
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
	return
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
	return
}

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event event.Event, payload interface{}) {
	m.Called(event, payload)
	return
}

type MockUserConfigRepository struct {
	mock.Mock
}

func (m *MockUserConfigRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockUserConfigRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

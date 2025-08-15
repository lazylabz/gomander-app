package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	"gomander/internal/event"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

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

func TestApp_GetUserConfig(t *testing.T) {
	mockRepository := new(MockRepository)

	a := &App{
		userConfigRepository: mockRepository,
	}

	expectedResult := domain.Config{
		LastOpenedProjectId: "test-project-id",
		EnvironmentPaths: []domain.EnvironmentPath{
			{
				Id:   "test-env-path-id",
				Path: "test/path",
			},
		},
	}
	mockRepository.On("GetOrCreate").Return(&expectedResult, nil)

	config, err := a.GetUserConfig()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, *config)
}

func TestApp_SaveUserConfig(t *testing.T) {
	t.Parallel()
	t.Run("Should save user configuration successfully", func(t *testing.T) {
		mockRepository := new(MockRepository)
		mockEventEmitter := new(MockEventEmitter)
		mockLogger := new(MockLogger)

		a := &App{
			userConfigRepository: mockRepository,
			eventEmitter:         mockEventEmitter,
			logger:               mockLogger,
		}

		newUserConfig := domain.Config{
			LastOpenedProjectId: "new-project-id",
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "new-env-path-id",
					Path: "new/path",
				},
			},
		}

		mockRepository.On("Update", &newUserConfig).Return(nil)

		mockLogger.On("Info", mock.Anything).Return(nil)

		mockEventEmitter.On("EmitEvent", event.SuccessNotification, mock.Anything).Return(nil)
		mockEventEmitter.On("EmitEvent", event.GetUserConfig, nil).Return(nil)

		err := a.SaveUserConfig(newUserConfig)
		assert.NoError(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
	})
	t.Run("Should fail to save user configuration", func(t *testing.T) {
		mockRepository := new(MockRepository)
		mockEventEmitter := new(MockEventEmitter)
		mockLogger := new(MockLogger)

		a := &App{
			userConfigRepository: mockRepository,
			eventEmitter:         mockEventEmitter,
			logger:               mockLogger,
		}

		newUserConfig := domain.Config{
			LastOpenedProjectId: "new-project-id",
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "new-env-path-id",
					Path: "new/path",
				},
			},
		}

		mockRepository.On("Update", &newUserConfig).Return(errors.New("failed to save user configuration"))

		mockLogger.On("Error", mock.Anything).Return(nil)

		mockEventEmitter.On("EmitEvent", event.ErrorNotification, mock.Anything).Return(nil)

		err := a.SaveUserConfig(newUserConfig)
		assert.Error(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
	})
}

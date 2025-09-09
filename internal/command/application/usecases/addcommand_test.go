package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/command/domain/test"
	configdomain "gomander/internal/config/domain"
	test2 "gomander/internal/config/domain/test"
)

func TestDefaultAddCommand_Execute(t *testing.T) {
	t.Run("Should add the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test2.MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		existingCommand := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()

		newCommandBuilder := test.NewCommandBuilder().
			WithProjectId(projectId)

		parameterCommand := newCommandBuilder.Build()

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			existingCommand,
		}, nil)
		expectedCommandCall := newCommandBuilder.WithPosition(1).Build()
		mockCommandRepository.On("Create", &expectedCommandCall).Return(nil)

		// Act
		err := sut.Execute(parameterCommand)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if fails to get the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test2.MockConfigRepository)

		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)
		newCommandBuilder := test.NewCommandBuilder()
		parameterCommand := newCommandBuilder.Build()
		expectedErr := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedErr)

		// Act
		err := sut.Execute(parameterCommand)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if fails to get all commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test2.MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		newCommandBuilder := test.NewCommandBuilder().
			WithProjectId(projectId)

		parameterCommand := newCommandBuilder.Build()

		mockCommandRepository.On("GetAll", projectId).Return(make([]commanddomain.Command, 0), errors.New("failed to get commands"))

		// Act
		err := sut.Execute(parameterCommand)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if fails to create commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test2.MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		existingCommand := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()

		newCommandBuilder := test.NewCommandBuilder().
			WithProjectId(projectId)

		parameterCommand := newCommandBuilder.Build()

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			existingCommand,
		}, nil)
		expectedCommandCall := newCommandBuilder.WithPosition(1).Build()
		mockCommandRepository.On("Create", &expectedCommandCall).Return(errors.New("failed to create command"))

		// Act
		err := sut.Execute(parameterCommand)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})
}

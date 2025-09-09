package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	commanddomain "gomander/internal/command/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/testutils"
)

func TestDefaultAddCommand_Execute(t *testing.T) {
	t.Run("Should add the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		existingCommandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			commandDataToDomain(existingCommandData),
		}, nil)
		expectedCommandCall := commandDataToDomain(newCommandData.WithPosition(1).Data())
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
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)
		newCommandData := testutils.NewCommand()
		parameterCommand := commandDataToDomain(newCommandData.Data())
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
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

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
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewAddCommand(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		existingCommandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			commandDataToDomain(existingCommandData),
		}, nil)
		expectedCommandCall := commandDataToDomain(newCommandData.WithPosition(1).Data())
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

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

func TestDefaultGetCommands_Execute(t *testing.T) {
	t.Run("Should return the commands provided by the repository", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewGetCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		command1Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		command2Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(1).
			Data()

		expectedCommandGroup := []commanddomain.Command{
			commandDataToDomain(command1Data),
			commandDataToDomain(command2Data),
		}

		mockCommandRepository.On("GetAll", projectId).Return(expectedCommandGroup, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, got, expectedCommandGroup)

		mock.AssertExpectationsForObjects(t, mockCommandRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if fails to get the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		sut := usecases.NewGetCommands(mockUserConfigRepository, mockCommandRepository)
		expectedErr := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedErr)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, got)
		mock.AssertExpectationsForObjects(t, mockCommandRepository, mockUserConfigRepository)
	})
}

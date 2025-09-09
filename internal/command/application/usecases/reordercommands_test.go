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
)

func TestDefaultReorderCommands_Execute(t *testing.T) {
	t.Run("Should reorder commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(0)
		cmd2 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(1)
		cmd3 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Build().Id, cmd1.Build().Id, cmd2.Build().Id}

		cm1WithUpdatedPosition := cmd1.WithPosition(1).Build()
		cm2WithUpdatedPosition := cmd2.WithPosition(2).Build()
		cm3WithUpdatedPosition := cmd3.WithPosition(0).Build()

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				cmd1.Build(),
				cmd2.Build(),
				cmd3.Build(),
			}, nil,
		)

		mockCommandRepository.On("Update", &cm1WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm2WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm3WithUpdatedPosition).Return(nil)

		// Act
		err := sut.Execute(orderedIds)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if fails to get the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)
		expectedErr := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedErr)

		// Act
		err := sut.Execute([]string{})

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if fails to retrieve commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		mockCommandRepository.On("GetAll", projectId).Return(
			make([]commanddomain.Command, 0),
			errors.New("failed to get commands"))

		// Act
		err := sut.Execute([]string{})

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if fails to update commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(0)
		cmd2 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(1)
		cmd3 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Build().Id, cmd1.Build().Id, cmd2.Build().Id}

		cm1WithUpdatedPosition := cmd1.WithPosition(1).Build()
		cm2WithUpdatedPosition := cmd2.WithPosition(2).Build()
		cm3WithUpdatedPosition := cmd3.WithPosition(0).Build()

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				cmd1.Build(),
				cmd2.Build(),
				cmd3.Build(),
			}, nil,
		)

		mockCommandRepository.On("Update", &cm1WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm2WithUpdatedPosition).Return(errors.New("failed to update command"))
		mockCommandRepository.On("Update", &cm3WithUpdatedPosition).Return(nil)

		// Act
		err := sut.Execute(orderedIds)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})
}

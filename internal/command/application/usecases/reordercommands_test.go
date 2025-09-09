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

func TestDefaultReorderCommands_Execute(t *testing.T) {
	t.Run("Should reorder commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := testutils.NewCommand().WithProjectId(projectId).WithPosition(0)
		cmd2 := testutils.NewCommand().WithProjectId(projectId).WithPosition(1)
		cmd3 := testutils.NewCommand().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Data().Id, cmd1.Data().Id, cmd2.Data().Id}

		cm1WithUpdatedPosition := commandDataToDomain(cmd1.WithPosition(1).Data())
		cm2WithUpdatedPosition := commandDataToDomain(cmd2.WithPosition(2).Data())
		cm3WithUpdatedPosition := commandDataToDomain(cmd3.WithPosition(0).Data())

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				commandDataToDomain(cmd1.Data()),
				commandDataToDomain(cmd2.Data()),
				commandDataToDomain(cmd3.Data()),
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
		mockCommandRepository := new(MockCommandRepository)
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
		mockCommandRepository := new(MockCommandRepository)
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
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewReorderCommands(mockUserConfigRepository, mockCommandRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := testutils.NewCommand().WithProjectId(projectId).WithPosition(0)
		cmd2 := testutils.NewCommand().WithProjectId(projectId).WithPosition(1)
		cmd3 := testutils.NewCommand().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Data().Id, cmd1.Data().Id, cmd2.Data().Id}

		cm1WithUpdatedPosition := commandDataToDomain(cmd1.WithPosition(1).Data())
		cm2WithUpdatedPosition := commandDataToDomain(cmd2.WithPosition(2).Data())
		cm3WithUpdatedPosition := commandDataToDomain(cmd3.WithPosition(0).Data())

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				commandDataToDomain(cmd1.Data()),
				commandDataToDomain(cmd2.Data()),
				commandDataToDomain(cmd3.Data()),
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

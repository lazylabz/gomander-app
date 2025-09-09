package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/testutils"
)

func TestDefaultReorderCommandGroups_Execute(t *testing.T) {
	t.Run("Should reorder command groups based on the provided IDs", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewReorderCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData1 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroupData2 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroups := []domain.CommandGroup{
			commandGroupDataToDomain(commandGroupData1.Data()),
			commandGroupDataToDomain(commandGroupData2.Data()),
		}

		newOrder := []string{commandGroupData2.Data().Id, commandGroupData1.Data().Id}

		mockCommandGroupRepository.On("GetAll", projectId).Return(commandGroups, nil)

		expectedCommandGroup2Call := commandGroupDataToDomain(commandGroupData2.WithPosition(0).Data())
		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData1.WithPosition(1).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup2Call).Return(nil).Once()
		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()

		// Act
		err := sut.Execute(newOrder)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve user config", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		sut := usecases.NewReorderCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)
		expectedError := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedError)

		newOrder := []string{"group1", "group2"}

		// Act
		err := sut.Execute(newOrder)

		// Assert
		assert.ErrorIs(t, err, expectedError)
		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewReorderCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		newOrder := []string{"group1", "group2"}

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))

		// Act
		err := sut.Execute(newOrder)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if failing to update a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewReorderCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData1 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroups := []domain.CommandGroup{
			commandGroupDataToDomain(commandGroupData1.Data()),
		}

		newOrder := []string{commandGroupData1.Data().Id}

		mockCommandGroupRepository.On("GetAll", projectId).Return(commandGroups, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData1.WithPosition(0).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(errors.New("failed to update command group")).Once()

		// Act
		err := sut.Execute(newOrder)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
}

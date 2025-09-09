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

func TestDefaultCreateCommandGroup_Execute(t *testing.T) {
	t.Run("Should create a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		otherCommandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			WithPosition(0).
			WithCommands(testutils.
				NewCommand().
				WithProjectId(projectId).
				Data(),
			)

		otherCommandGroup := commandGroupDataToDomain(otherCommandGroupData)
		paramCommandGroup := commandGroupDataToDomain(commandGroupData.Data())
		expectedCommandGroupCall := commandGroupDataToDomain(commandGroupData.WithPosition(1).Data())

		mockCommandGroupRepository.On("GetAll", projectId).Return([]domain.CommandGroup{otherCommandGroup}, nil)
		mockCommandGroupRepository.On("Create", &expectedCommandGroupCall).Return(nil)

		// Act
		err := sut.Execute(&paramCommandGroup)

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

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		expectedError := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedError)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId("project1").
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to save the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), nil)
		mockCommandGroupRepository.On("Create", &paramCommandGroup).Return(errors.New("failed to create command group"))

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
}
